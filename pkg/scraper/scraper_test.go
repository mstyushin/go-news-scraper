package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
	"go-news-scraper/pkg/storage/pg"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbURL = "postgres://postgres@localhost:5432/news?sslmode=disable"

var s *pg.DB
var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	s, err = pg.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}

	testMux := mux.NewRouter()
	testMux.HandleFunc("/rss", getRss).Methods(http.MethodGet, http.MethodOptions)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", 12345),
		Handler: testMux,
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.WithValue(context.Background(), fmt.Sprintf(":%v", 12345), l.Addr().String())
			return ctx
		},
	}
	go func(s *http.Server) {
		if err := s.ListenAndServe(); err != nil {
			os.Exit(1)
		}
	}(httpServer)

	os.Exit(m.Run())
}

func getRss(w http.ResponseWriter, r *http.Request) {
	rawXml := `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Golang Weekly</title>
    <description>A weekly newsletter about the Go programming language</description>
    <link>https://golangweekly.com/</link>
    <item>
      <title>Channels that channel channels?</title>
      <link>https://golangweekly.com/issues/520</link>
      <description>
&lt;table border=0 cellpadding=0 cellspacing=0 align="center" border="0"&gt;
  &lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;
  &lt;div&gt;    
    &lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;
&lt;td align="left" style="padding-left: 4px; font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;&lt;p&gt;#​520 — August 27, 2024&lt;/p&gt;&lt;/td&gt;
&lt;td align="right" style="padding-right: 4px; font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;&lt;p&gt;&lt;a href="https://golangweekly.com/link/158911/rss" style=" color: #0099b4;"&gt;Unsub&lt;/a&gt;  |  &lt;a href="https://golangweekly.com/link/158912/rss" style=" color: #0099b4;"&gt;Web Version&lt;/a&gt;&lt;/p&gt;&lt;/td&gt;
&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style=" font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
    
    &lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0; padding-right: 12px;  padding-left: 12px;"&gt;&lt;p&gt;Go Weekly&lt;/p&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;
  &lt;a href="https://golangweekly.com/link/158914/rss" style=" color: #0099b4;"&gt;&lt;img src="https://res.cloudinary.com/cpress/image/upload/w_1280,e_sharpen:60,q_auto/gwoxp7krr7e4upv03u37.jpg" width="640" style="    line-height: 100%;       "&gt;&lt;/a&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158914/rss" title="go.dev" style=" color: #0099b4;    font-size: 1.1em; line-height: 1.4em;"&gt;Range Over Function Types in Go 1.23&lt;/a&gt;&lt;/span&gt; — &lt;a href="https://golangweekly.com/link/158915/rss" style=" color: #0099b4;   "&gt;Go 1.23&lt;/a&gt; was released just two weeks ago with one of the headline features being improvements to iteration, principally &lt;code&gt;for/range&lt;/code&gt; support over function types. There’s a lot more to it than that, though, and Ian rounds up everything, complete with examples and guidance.&lt;/p&gt;
  &lt;p&gt;Ian Lance Taylor &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158916/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;TinyGo 0.33.0: The Go Compiler for 'Small Places'&lt;/a&gt;&lt;/span&gt; — &lt;a href="https://golangweekly.com/link/158917/rss" style=" color: #0099b4;   "&gt;TinyGo&lt;/a&gt; is a fantastic LLVM-based Go implementation targeting use cases like microcontrollers and WebAssembly. v0.33.0 brings it up to Go 1.23 standards, adds WASI preview 2 support, and more.&lt;/p&gt;
  &lt;p&gt;TinyGo Team &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  &lt;a href="https://golangweekly.com/link/158913/rss" style=" color: #0099b4;   "&gt;&lt;img src="https://copm.s3.amazonaws.com/1662bfe0.png" width="127" height="110" style="padding-top: 12px; padding-left: 12px;     line-height: 100%;    "&gt;&lt;/a&gt;
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158913/rss" title="workos.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;WorkOS: The Modern Identity Platform for B2B SaaS&lt;/a&gt;&lt;/span&gt; — WorkOS is a modern identity platform for B2B SaaS, offering flexible and easy-to-use APIs to integrate SSO, SCIM, and RBAC in minutes instead of months. It's trusted by hundreds of high-growth startups such as Perplexity, Vercel, Drata, and Webflow.&lt;/p&gt;
  &lt;p&gt;WorkOS &lt;span style="text-transform: uppercase; margin-left: 4px; font-size: 0.9em;   color: #885 !important; padding-top: 1px; padding-right: 4px;  padding-left: 4px;            "&gt;sponsor&lt;/span&gt;&lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158918/rss" title="www.dolthub.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;The 4-&lt;code&gt;chan&lt;/code&gt; Go Programmer (or Sending Channels Over Channels)&lt;/a&gt;&lt;/span&gt; — No, not &lt;a href="https://golangweekly.com/link/158919/rss" style=" color: #0099b4;   "&gt;&lt;em&gt;that&lt;/em&gt;&lt;/a&gt; '4chan.' Zach looks into the idea of using channels to pass &lt;em&gt;other&lt;/em&gt; channels around, and takes it to the extreme of creating channels that channel channels that channel channels that channel channels. &lt;em&gt;I need to go for a lie down..&lt;/em&gt;&lt;/p&gt;
  &lt;p&gt;Zach Musgrave &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158920/rss" title="victoriametrics.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;&lt;code&gt;sync.Pool&lt;/code&gt; and the Mechanics Behind It&lt;/a&gt;&lt;/span&gt; — What &lt;code&gt;sync.Pool&lt;/code&gt; is all about, how it’s used, what’s going on under the hood, and everything else you might want to know, complete with helpful illustrations.&lt;/p&gt;
  &lt;p&gt;Phuong Le &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
&lt;p&gt;📄 &lt;a href="https://golangweekly.com/link/158921/rss" style=" color: #0099b4; font-weight: 500 !important;"&gt;Designing a Robust Integration Test Suite for Convoy’s Data Plane with TestContainers&lt;/a&gt;  &lt;cite&gt;Oluwalana and Mekiliuwa (Convoy)&lt;/cite&gt;&lt;/p&gt;
&lt;p&gt;📄 &lt;a href="https://golangweekly.com/link/158922/rss" style=" color: #0099b4; font-weight: 500 !important;"&gt;'Go is My Hammer, and Everything is a Nail'&lt;/a&gt; – You can use Go for everything. &lt;cite&gt;Markus Wüstenberg&lt;/cite&gt;&lt;/p&gt;
&lt;p&gt;📄 &lt;a href="https://golangweekly.com/link/158923/rss" style=" color: #0099b4; font-weight: 500 !important;"&gt;How Go's Testing Harness Works&lt;/a&gt; – What happens when you run &lt;code&gt;go test&lt;/code&gt;? &lt;cite&gt;Matt Proud&lt;/cite&gt;&lt;/p&gt;
&lt;p&gt;📄 &lt;a href="https://golangweekly.com/link/158924/rss" style=" color: #0099b4; font-weight: 500 !important;"&gt;Using Functional Options Instead of Method Chaining&lt;/a&gt;  &lt;cite&gt;Jon Calhoun&lt;/cite&gt;&lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0; padding-right: 0;  padding-left: 0;"&gt;&lt;p&gt;🛠 Code &amp;amp; Tools&lt;/p&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;
  &lt;a href="https://golangweekly.com/link/158925/rss" style=" color: #0099b4;"&gt;&lt;img src="https://res.cloudinary.com/cpress/image/upload/w_1280,e_sharpen:60,q_auto/shl0c83idvrdnu9diaun.jpg" width="640" style="    line-height: 100%;         "&gt;&lt;/a&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158925/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;sqlite-vec: A Vector Search Extension for SQLite&lt;/a&gt;&lt;/span&gt; — If using a dedicated vector storage database is beyond your immediate requirements, you can now use SQLite for the task. While this is an SQLite extension, it’s easy to use with numerous languages, &lt;a href="https://golangweekly.com/link/158926/rss" style=" color: #0099b4;   "&gt;including Go, as shown here.&lt;/a&gt;&lt;/p&gt;
  &lt;p&gt;Alex Garcia &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158927/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;moq: An Interface Mocking Tool for &lt;code&gt;go generate&lt;/code&gt;&lt;/a&gt;&lt;/span&gt; — A tool that generates a struct from any interface. The struct can be used in test code as a mock of the interface. Now supports imported type aliases.&lt;/p&gt;
  &lt;p&gt;Mat Ryer &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158928/rss" title="monday.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;Streamline Your Product Delivery with monday dev&lt;/a&gt;&lt;/span&gt; — From ideation to launch, monday dev makes product delivery faster and simpler, all in one place.&lt;/p&gt;
  &lt;p&gt;monday dev &lt;span style="text-transform: uppercase; margin-left: 4px; font-size: 0.9em;   color: #885 !important; padding-top: 1px; padding-right: 4px;  padding-left: 4px;            "&gt;sponsor&lt;/span&gt;&lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158929/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;Permify 1.0: Open Source Authorization as a Service&lt;/a&gt;&lt;/span&gt; — A long standing Go-powered system for building fine-grained authorization systems, inspired by &lt;a href="https://golangweekly.com/link/158930/rss" style=" color: #0099b4;   "&gt;Google’s Zanzibar.&lt;/a&gt; Get started with the &lt;a href="https://golangweekly.com/link/158931/rss" style=" color: #0099b4;   "&gt;intro guide.&lt;/a&gt;&lt;/p&gt;
  &lt;p&gt;Permify &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158932/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;cast 1.7: Safe and Easy Casting from One Type to Another&lt;/a&gt;&lt;/span&gt; — It’s as easy as using &lt;code&gt;ToString&lt;/code&gt;, &lt;code&gt;ToInt&lt;/code&gt;, &lt;code&gt;ToTime&lt;/code&gt;, and more.&lt;/p&gt;
  &lt;p&gt;Steve Francia &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158933/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;go-github v64.0: A Go Client for the GitHub v3 API&lt;/a&gt;&lt;/span&gt; — For the REST API. For the v4 GraphQL API &lt;a href="https://golangweekly.com/link/158934/rss" style=" color: #0099b4;   "&gt;use this instead.&lt;/a&gt;&lt;/p&gt;
  &lt;p&gt;Google &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158935/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;Dbmate: A Lightweight, Framework-Agnostic Database Migration Tool&lt;/a&gt;&lt;/span&gt; — Written in Go but can be used alongside database-using apps written in any language. Supports MySQL, Postgres, SQLite, ClickHouse, BigQuery, and Spanner.&lt;/p&gt;
  &lt;p&gt;Adrian Macneil &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 0px;  padding-left: 0px;"&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style=" font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
	&lt;p&gt;📰 Classifieds&lt;/p&gt;
  &lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158936/rss" style=" color: #0099b4; font-weight: 500 !important;"&gt;Dragonfly (25k GitHub stars)&lt;/a&gt; is a modern Redis replacement. Organizations that switch to Dragonfly can reduce infrastructure costs by 80%.&lt;/p&gt;
 
&lt;p&gt;Boost your PostgreSQL skills with &lt;a href="https://golangweekly.com/link/158937/rss" style=" color: #0099b4; font-weight: 500 !important;"&gt;Redgate’s 101 webinar series&lt;/a&gt; of easy-to-follow, expert hosted sessions. It’s PostgreSQL, simplified.&lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style=" font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style=" font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
&lt;ul&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158938/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;go-rest-api-example&lt;/a&gt; – Template for an enterprise-ready REST API microservice.&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158939/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Glow 2.0&lt;/a&gt; – Markdown renderer for the terminal.&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158940/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Dynamo 2.2&lt;/a&gt; – Expressive &lt;em&gt;Amazon DynamoDB&lt;/em&gt; client library.&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;💄 &lt;a href="https://golangweekly.com/link/158941/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Charm Lip Gloss 0.13&lt;/a&gt; – Style definitions for rendering nice terminal layouts. It can now render trees too!&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158942/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;GoWrap 1.4&lt;/a&gt; – Tool to generate decorators for Go interfaces.&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158943/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Buf 1.38&lt;/a&gt; – CLI tool for working with Protocol Buffers.&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158944/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;c-for-go 1.3&lt;/a&gt; – Automatic C-Go bindings generator.&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158945/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Go OpenAI 1.29&lt;/a&gt; – An OpenAI API wrapper library.&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158946/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Bubbles 0.19&lt;/a&gt; – TUI components for &lt;a href="https://golangweekly.com/link/158947/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Bubble Tea.&lt;/a&gt;&lt;/p&gt;
&lt;/li&gt;
&lt;li&gt;
&lt;p&gt;&lt;a href="https://golangweekly.com/link/158948/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;netlink 1.3&lt;/a&gt; – Simple &lt;a href="https://golangweekly.com/link/158949/rss" style=" color: #0099b4; font-weight: 500 !important;   "&gt;Netlink&lt;/a&gt; library for Go.&lt;/p&gt;
&lt;/li&gt;
&lt;/ul&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0; padding-right: 0;  padding-left: 0;"&gt;&lt;p&gt;🎁 And one for fun..&lt;/p&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;
  &lt;a href="https://golangweekly.com/link/158950/rss" style=" color: #0099b4;"&gt;&lt;img src="https://res.cloudinary.com/cpress/image/upload/w_1280,e_sharpen:60,q_auto/nbd25nbblmn3juzd86yl.jpg" width="640" style="        line-height: 100%;     "&gt;&lt;/a&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;

&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style="font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em;  padding-top: 0px; padding-right: 15px;  padding-left: 15px;"&gt;
  
  &lt;p&gt;&lt;span style="font-weight: 600; font-size: 1.2em !important; color: #000;"&gt;&lt;a href="https://golangweekly.com/link/158950/rss" title="github.com" style=" color: #0099b4;    font-size: 1.05em;"&gt;Tetrigo: TUI-Powered Tetris Written in Go&lt;/a&gt;&lt;/span&gt; — Whether you want to just play Tetris, create your own Tetris game, or perhaps implement your own quirky Tetris game mode using Go, this implementation is for you. It’s well structured and uses Charm’s &lt;a href="https://golangweekly.com/link/158947/rss" style=" color: #0099b4;   "&gt;Bubble Tea&lt;/a&gt; behind the scenes.&lt;/p&gt;
  &lt;p&gt;Broderick Westrope &lt;/p&gt;
&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;table border=0 cellpadding=0 cellspacing=0 border=0 cellpadding=0 cellspacing=0&gt;&lt;tr&gt;&lt;td style=" font-family: -apple-system,BlinkMacSystemFont,Helvetica,sans-serif; font-size: 15px; line-height: 1.48em; "&gt;&lt;/td&gt;&lt;/tr&gt;&lt;/table&gt;
&lt;/div&gt;
  &lt;/td&gt;&lt;/tr&gt;
&lt;/table&gt;

&lt;img src="https://golangweekly.com/open/520/rss" width="1" height="1" /&gt;</description>
      <pubDate>Tue, 27 Aug 2024 00:00:00 +0000</pubDate>
      <guid>https://golangweekly.com/issues/520</guid>
    </item>
  </channel>
</rss>
`
	bytes, _ := xml.Marshal(rawXml)
	w.Write(bytes)

}

// недоделанный тест, к сожалению не успеваю отладить
func TestScraper_scrape(t *testing.T) {
	scraper := New(s, 10*time.Minute)
	scraper.Run(context.Background())
	time.Sleep(1 * time.Second)
}
