package ttt

/*
 *
 * http://localhost:8080/api/v1/twitter/tokenurl/redirect?cb=http%3A%2F%2Flocalhost%3A8080%2Ftest
 *
 * http://localhost:8080/api/v1/facebook/tokenurl/redirect?cb=http%3A%2F%2Flocalhost%3A8080%2Ftest
 *
 * http://localhost:8080/api/v1/facebook/tokenurl/callback
 */
import (
	"net/http"

	"github.com/firefirestyle/go.minioauth/facebook"
	"github.com/firefirestyle/go.minioauth/twitter"
	//
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	//
)

const (
	UrlTwitterTokenUrlRedirect  = "/api/v1/twitter/tokenurl/redirect"
	UrlTwitterTokenCallback     = "/api/v1/twitter/tokenurl/callback"
	UrlFacebookTokenUrlRedirect = "/api/v1/facebook/tokenurl/redirect"
	UrlFacebookTokenCallback    = "/api/v1/facebook/tokenurl/callback"
)

var twitterHandlerObj *twitter.TwitterHandler = nil

var facebookHandlerObj *facebook.FacebookHandler = nil

func GetTwitterHandlerObj(ctx context.Context) *twitter.TwitterHandler {
	if twitterHandlerObj == nil {
		twitterHandlerObj = twitter.NewTwitterHandler( //
			twitter.TwitterOAuthConfig{
				ConsumerKey:       TwitterConsumerKey,
				ConsumerSecret:    TwitterConsumerSecret,
				AccessToken:       TwitterAccessToken,
				AccessTokenSecret: TwitterAccessTokenSecret,
				SecretSign:        "abc",
				CallbackUrl:       "http://" + appengine.DefaultVersionHostname(ctx) + "" + UrlTwitterTokenCallback,
			},
			twitter.TwitterHundlerOnEvent{
				OnRequest: func(http.ResponseWriter, *http.Request, *twitter.TwitterHandler) (map[string]string, error) {
					return map[string]string{"test": "abcdef"}, nil
				},
				OnFoundUser: func(w http.ResponseWriter, r *http.Request, h *twitter.TwitterHandler, s *twitter.SendAccessTokenResult) map[string]string {
					return map[string]string{ //
						"test":       r.URL.Query().Get("test"), //
						"userId":     s.GetUserID(),             //
						"screenName": s.GetScreenName(),
						"token":      s.GetOAuthToken(),
						"secret":     s.GetOAuthTokenSecret(),
					}
				},
			})
	}
	return twitterHandlerObj
}

func GetFacebookHandlerObj(ctx context.Context) *facebook.FacebookHandler {
	v := appengine.DefaultVersionHostname(ctx)
	if v == "127.0.0.1:8080" {
		v = "localhost:8080"
	}
	if facebookHandlerObj == nil {
		facebookHandlerObj = facebook.NewFacebookHandler( //
			facebook.FacebookOAuthConfig{
				ConfigFacebookAppSecret: ConfigFacebookAppSecret,
				ConfigFacebookAppId:     ConfigFacebookAppId,
				CallbackUrl:             "http://" + v + "" + UrlFacebookTokenCallback,
				SecretSign:              "abc",
			},
			facebook.FacebookHundlerOnEvent{
				OnRequest: func(http.ResponseWriter, *http.Request, *facebook.FacebookHandler) (map[string]string, error) {
					return map[string]string{"test": "abcdef"}, nil
				},
				OnFoundUser: func(w http.ResponseWriter, r *http.Request, h *facebook.FacebookHandler, s *facebook.GetMeResponse, t *facebook.AccessTokenResponse) map[string]string {
					return map[string]string{
						"test":       r.URL.Query().Get("test"), //
						"userId":     s.Id,                      //
						"screenName": s.Name,
						"token":      t.AccessToken,
					}
				},
			})
	}
	return facebookHandlerObj
}
func init() {
	initApi()
	initHomepage()
}

func initHomepage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FireFireStyle!!"))
	})
}

func initApi() {
	// twitter
	http.HandleFunc(UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).HandleLoginEntry(w, r)
	})
	http.HandleFunc(UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).HandleLoginExit(w, r)
	})

	http.HandleFunc(UrlFacebookTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetFacebookHandlerObj(appengine.NewContext(r)).HandleLoginEntry(w, r)
	})
	http.HandleFunc(UrlFacebookTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetFacebookHandlerObj(appengine.NewContext(r)).HandleLoginExit(w, r)
	})
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
