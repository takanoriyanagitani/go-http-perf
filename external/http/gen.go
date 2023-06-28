package hgen

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	util "github.com/takanoriyanagitani/go-http-perf/util"

	ext "github.com/takanoriyanagitani/go-http-perf/external"
)

type Url struct {
	raw string
}

func UrlNewFromRawString(raw string) Url { return Url{raw} }

func (u Url) WithSuffix(suffix string) Url {
	var neo string = u.raw + suffix
	return Url{raw: neo}
}

func (u Url) Get(ctx context.Context) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, u.raw, nil)
}

type UrlConverter func(original Url) (converted Url)

type UrlSuffixGenerator[T any] func(T) (suffix string)

func (s UrlSuffixGenerator[T]) ToUrlConverter(gen func() T) UrlConverter {
	return func(original Url) (converted Url) {
		var input T = gen()
		var suffix string = s(input)
		return original.WithSuffix(suffix)
	}
}

type TimeSerializer func(time.Time) (serialized string)

var TimeSerializerMicro TimeSerializer = util.Compose(
	func(t time.Time) (us int64) { return t.UnixMicro() },
	func(us int64) string { return strconv.FormatInt(us, 10) },
)

type TimeGenerator func() time.Time

func (s TimeSerializer) ToSuffixGen() UrlSuffixGenerator[time.Time] {
	return func(t time.Time) (suffix string) {
		var serialized string = s(t)
		return "?now=" + serialized
	}
}
func (s TimeSerializer) ToUrlConverter(gen TimeGenerator) UrlConverter {
	var sgen UrlSuffixGenerator[time.Time] = s.ToSuffixGen()
	return sgen.ToUrlConverter(gen)
}

var TimeGeneratorNow TimeGenerator = time.Now

var UrlConverterDefault UrlConverter = TimeSerializerMicro.ToUrlConverter(TimeGeneratorNow)

type SerializedRequestGeneratorBuilder struct {
	url Url
	cnv UrlConverter
	cli *http.Client
}

func (b SerializedRequestGeneratorBuilder) WithUrl(u Url) SerializedRequestGeneratorBuilder {
	return SerializedRequestGeneratorBuilder{
		url: u,
		cnv: b.cnv,
		cli: b.cli,
	}
}

func (b SerializedRequestGeneratorBuilder) WithCnv(c UrlConverter) SerializedRequestGeneratorBuilder {
	return SerializedRequestGeneratorBuilder{
		url: b.url,
		cnv: c,
		cli: b.cli,
	}
}

func (b SerializedRequestGeneratorBuilder) WithCli(c *http.Client) SerializedRequestGeneratorBuilder {
	return SerializedRequestGeneratorBuilder{
		url: b.url,
		cnv: b.cnv,
		cli: c,
	}
}

var DefaultBuilder SerializedRequestGeneratorBuilder = SerializedRequestGeneratorBuilder{}.
	WithUrl(UrlNewFromRawString("http://localhost")).
	WithCnv(UrlConverterDefault).
	WithCli(http.DefaultClient)

func res2bytes(res *http.Response) ([]byte, error) {
	defer func() {
		_ = res.Body.Close()
	}()
	return io.ReadAll(res.Body)
}

func (s SerializedRequestGeneratorBuilder) Build() ext.SerializedRequestGenerator {
	return func(ctx context.Context) (serialized []byte, e error) {
		var converted Url = s.cnv(s.url)
		var ctx2res func(context.Context) (*http.Response, error) = util.ComposeErr(
			converted.Get,
			s.cli.Do,
		)
		var ctx2bytes func(context.Context) ([]byte, error) = util.ComposeErr(
			ctx2res,
			res2bytes,
		)
		return ctx2bytes(ctx)
	}
}
