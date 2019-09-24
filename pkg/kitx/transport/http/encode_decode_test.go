package http

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/moogar0880/problems"
)

func TestNopResponseEncoder(t *testing.T) {
	handler := kithttp.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return "response", nil },
		kithttp.NopRequestDecoder,
		NopResponseEncoder,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("unexpected status code\nexpected: %d\nactual:   %d", want, have)
	}
}

func TestStatusCodeResponseEncoder(t *testing.T) {
	handler := kithttp.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return "response", nil },
		kithttp.NopRequestDecoder,
		StatusCodeResponseEncoder(http.StatusNoContent),
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := http.StatusNoContent, resp.StatusCode; want != have {
		t.Errorf("unexpected status code\nexpected: %d\nactual:   %d", want, have)
	}
}

func TestJSONResponseEncoder(t *testing.T) {
	handler := kithttp.NewServer(
		func(context.Context, interface{}) (interface{}, error) {
			return struct {
				Foo string `json:"foo"`
			}{Foo: "bar"}, nil
		},
		kithttp.NopRequestDecoder,
		JSONResponseEncoder,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("unexpected status code\nexpected: %d\nactual:   %d", want, have)
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := `{"foo":"bar"}`, strings.TrimSpace(string(buf)); want != have {
		t.Errorf("unexpected body\nexpected: %s\nactual:   %s", want, have)
	}
}

type failer struct {
	err error
}

func (f failer) Failed() error {
	return f.err
}

func TestErrorResponseEncoder(t *testing.T) {
	t.Parallel()

	t.Run("response", func(t *testing.T) {
		handler := kithttp.NewServer(
			func(context.Context, interface{}) (interface{}, error) {
				return struct {
					Foo string `json:"foo"`
				}{Foo: "bar"}, nil
			},
			kithttp.NopRequestDecoder,
			ErrorResponseEncoder(JSONResponseEncoder, func(i context.Context, w http.ResponseWriter, e error) error {
				problem := problems.NewDetailedProblem(http.StatusBadRequest, e.Error())

				w.Header().Set("Content-Type", problems.ProblemMediaType)
				w.WriteHeader(problem.Status)

				return json.NewEncoder(w).Encode(problem)
			}),
		)

		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}

		if want, have := http.StatusOK, resp.StatusCode; want != have {
			t.Errorf("unexpected status code\nexpected: %d\nactual:   %d", want, have)
		}

		buf, _ := ioutil.ReadAll(resp.Body)
		if want, have := `{"foo":"bar"}`, strings.TrimSpace(string(buf)); want != have {
			t.Errorf("unexpected body\nexpected: %s\nactual:   %s", want, have)
		}
	})

	t.Run("error", func(t *testing.T) {
		handler := kithttp.NewServer(
			func(context.Context, interface{}) (interface{}, error) {
				return failer{errors.New("error")}, nil
			},
			kithttp.NopRequestDecoder,
			ErrorResponseEncoder(JSONResponseEncoder, func(i context.Context, w http.ResponseWriter, e error) error {
				problem := problems.NewDetailedProblem(http.StatusBadRequest, e.Error())

				w.Header().Set("Content-Type", problems.ProblemMediaType)
				w.WriteHeader(problem.Status)

				return json.NewEncoder(w).Encode(problem)
			}),
		)

		server := httptest.NewServer(handler)
		defer server.Close()

		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}

		if want, have := http.StatusBadRequest, resp.StatusCode; want != have {
			t.Errorf("unexpected status code\nexpected: %d\nactual:   %d", want, have)
		}

		buf, _ := ioutil.ReadAll(resp.Body)
		if want, have := `{"type":"about:blank","title":"Bad Request","status":400,"detail":"error"}`, strings.TrimSpace(string(buf)); want != have {
			t.Errorf("unexpected body\nexpected: %s\nactual:   %s", want, have)
		}
	})
}

func TestProblemErrorEncoder(t *testing.T) {
	handler := kithttp.NewServer(
		func(context.Context, interface{}) (interface{}, error) {
			return nil, errors.New("error")
		},
		kithttp.NopRequestDecoder,
		JSONResponseEncoder,
		kithttp.ServerErrorEncoder(ProblemErrorEncoder),
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("unexpected status code\nexpected: %d\nactual:   %d", want, have)
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := `{"type":"about:blank","title":"Internal Server Error","status":500,"detail":"something went wrong"}`, strings.TrimSpace(string(buf)); want != have {
		t.Errorf("unexpected body\nexpected: %s\nactual:   %s", want, have)
	}
}
