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

const CollectionDefaultTemplate = `
<!DOCTYPE html>
<html lang="{{ .LanguageCode }}">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{ .Title }}</title>
  </head>
  <body>
    <header>
      <h1>{{ .Title }}</h1>
    </header>

    <main>
      <h2>Books</h2>
      <ul>
        {{ range .Books -}}
        <li>
          <a href="{{ .ID }}/index.html">
            <em class="collection-book-title">{{ .Title }}</em>
            by
            <em class="collection-book-author">{{ .Author.Name }}</em>
          </a>
        </li>
        {{- end }}
      </ul>
    </main>
  </body>
</html>
`

const BookDefaultTemplate = `
<!DOCTYPE html>
<html lang="{{ .LanguageCode }}">
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
      <img src="{{ .CoverPath }}" alt="cover image for {{ .Title }}" style="display:block;">
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
        <li>Genre: {{ .Genre }}</li>
        <li>Status: {{ .Status }}</li>
        {{ with .Mirrors }}
        <li>
          Mirrors:
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

      {{ with .Author -}}
      <article name="about the author">
        <h2>About {{ .Name }}</h2>
        {{ .Bio }}

        {{ with .Donate -}}
        <ul>
          {{ range . -}}
          <li>
            {{ if .NonLinkItem -}}
            {{ .Name }}: {{ .Link }}
            {{- else -}}
            <a href="{{ .Link }}">{{ .Name }}</a>
            {{- end }}
          </li>
          {{- end }}
        </ul>
        {{- end }}
      </article>
      {{- end }}

      <h2>Table of Contents</h2>
      <nav>
        <ol>
          {{ range .Chapters }}
          <li>
            <a href="{{ .SlugHTML }}">{{ .Title }}</a>
          </li>
          {{ end }}
        </ol>
      </nav>
    </main>

    <footer>
      {{ .Copyright }}. Licensed under {{ .License }}
    </footer>
  </body>
</html>
`

const ChapterDefaultTemplate = `
<!DOCTYPE html>
<html lang="{{ .Parent.LanguageCode }}">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>{{ .Title }}</title>
  </head>
  <body>
    <header>
      <h1>{{ .Title }}</h1>
      <ul>
        {{ if not .Published.IsZero -}}
        <li>Published {{ .Published.Format "January 2 2006 3:04 PM PST" }}</li>
        {{- end }}
        {{ if not .LastModified.IsZero -}}
        <li>Last Modified {{ .LastModified.Format "January 2 2006 3:04 PM PST" }}</li>
        {{- end }}
        <li>{{ .EstimatedReadingTime.Words }} Words</li>
        <li>{{ .EstimatedReadingTime.Text }}</li>
      </ul>
      <h2><a href="index.html">{{ .Parent.Title }}</a></h2>

      {{ with .Description -}}
      {{ .Description }}
      {{- end }}
    </header>

    <main>
      {{ .Content }}
    </main>

    <footer>
      <nav>
        {{ with .Prev -}}
        <a href="{{ .SlugHTML }}">Previous: {{ .Title }}</a>
        {{- end }}
        {{ with .Next -}}
        <a href="{{ .SlugHTML }}">Next: {{ .Title }}</a>
        {{- end }}
      </nav>
      {{ .Parent.Copyright }}. Licensed under {{ .Parent.License }}
    </footer>
  </body>
</html>
`
