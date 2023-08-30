package renderer

const RSSTemplate = `
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ .Title }}</title>
    <link>{{ .Collection.BaseURL }}/{{ .ID }}</link>
    <description>Recent content for {{ .Title }}</description>
    <generator>Bookgen -- github.com/JessebotX/bookgen</generator>
    {{ with .LanguageCode -}}
    <language>{{.}}</language>
    {{- end }}
    {{ with .Copyright -}}
    <copyright>{{.}}</copyright>
    {{- end }}
    {{ range .Chapters -}}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Collection.BaseURL }}/{{ .Parent.ID }}/{{ .SlugHTML }}</link>
      <pubDate>{{ .Published.Format "Mon, 02 Jan 2006 15:04:05 -0700" }}</pubDate>
      <guid>{{ .Collection.BaseURL }}/{{ .Parent.ID }}/{{ .SlugHTML }}</guid>
      <description>
        {{- with .Description }}
        {{ . }}
        {{- else -}}
        {{ .Title }}
        {{ end -}}
      </description>
    </item>
    {{- end }}
  </channel>
</rss>
`

const BookDefaultTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{ .Title }}</title>
  </head>
  <body>
    <header>
      <h1>{{ .Title }}</h1>
      <em>by {{ .Author.Name }}</em>
      <p>{{ .ShortDescription }}</p>
    </header>

    <main>
      {{ if not (eq .CoverPath "") }}
      <div>
        <img src="{{ .CoverPath }}" alt="cover image for {{ .Title }}">
      </div>
      {{ end }}

      <article name="about">
        <h2>About</h2>
        {{ .Blurb }}
      </article>

      <h2>Stats</h2>
      <ul>
        <li>
          <a href="{{ .ID }}.epub">Download EPUB (offline)</a>
        </li>
        <li>
          <a href="./rss.xml">Follow RSS feed for updates</a>
        </li>
        <li>{{ .Genre }}</li>
        <li>{{ .Status }}</li>
        {{ with .Mirrors }}
        <li>
          Mirrors
          <ul>
            {{ range . }}
            <li>
              <a href="{{ . }}">{{ . }}</a>
            </li>
            {{ end }}
          </ul>
        </li>
        {{ end }}
      </ul>

      <h2>Table of Contents</h2>
      <ol>
        {{ range .Chapters }}
        <li>
          <a href="{{ .SlugHTML }}">{{ .Title }}</a>
        </li>
        {{ end }}
      </ol>
    </main>

    <footer>
      {{ .Copyright }}. Licensed under {{ .License }}
    </footer>
  </body>
</html>
`
