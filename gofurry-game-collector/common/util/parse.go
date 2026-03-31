package util

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
)

// Steam BBCode 特殊标签替换
var bbReplacements = map[string]string{
	// 文本样式
	`(?i)\[b\](.*?)\[/b\]`:     `<strong>$1</strong>`,
	`(?i)\[i\](.*?)\[/i\]`:     `<em>$1</em>`,
	`(?i)\[u\](.*?)\[/u\]`:     `<u>$1</u>`,
	`(?i)\[s\](.*?)\[/s\]`:     `<s>$1</s>`,
	`(?i)\[h1\](.*?)\[/h1\]`:   `<h1>$1</h1>`,
	`(?i)\[h2\](.*?)\[/h2\]`:   `<h2>$1</h2>`,
	`(?i)\[h3\](.*?)\[/h3\]`:   `<h3>$1</h3>`,
	`(?i)\[p\](.*?)\[/p\]`:     `$1<br>`,
	`(?i)\[img\](.*?)\[/img\]`: `<img src="$1" />`,

	// 链接
	`(?i)\[url\](.*?)\[/url\]`:       `<a href="$1">$1</a>`,
	`(?i)\[url=(.*?)\](.*?)\[/url\]`: `<a href="$1">$2</a>`,
}

// 将 BBCode 递归解析成 HTML
func ParseBBCode(input string) string {
	// 先替换 Steam 图片前缀
	input = strings.ReplaceAll(input, "{STEAM_CLAN_IMAGE}", "https://clan.fastly.steamstatic.com/images/")

	// 处理 img/video/youtube
	input = regexp.MustCompile(`(?i)\[img\s+src="(.*?)"\]\s*\[/img\]`).ReplaceAllString(input, `<img src="$1" />`)
	input = regexp.MustCompile(`(?i)\[video\](.*?)\[/video\]`).ReplaceAllString(input, `<video src="$1" controls></video>`)
	input = regexp.MustCompile(`(?i)\[youtube\](.*?)\[/youtube\]`).ReplaceAllString(input, `<iframe src="https://www.youtube.com/embed/$1" frameborder="0" allowfullscreen></iframe>`)

	// 替换常规 BBCode 标签
	for pattern, replacement := range bbReplacements {
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllString(input, replacement)
	}
	// 解析第一次嵌套
	for pattern, replacement := range bbReplacements {
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllString(input, replacement)
	}
	// 解析第二次嵌套
	for pattern, replacement := range bbReplacements {
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllString(input, replacement)
	}

	// 递归解析列表
	input = parseLists(input)

	// 将换行转换为 <br>
	input = strings.ReplaceAll(input, "\n", "<br>")

	return input
}

// 递归解析 [list] 和 [olist] 以及 [*] 项
func parseLists(input string) string {
	listPattern := regexp.MustCompile(`(?is)\[list\](.*?)\[/list\]`)
	olistPattern := regexp.MustCompile(`(?is)\[olist\](.*?)\[/olist\]`)

	// 处理无序列表
	input = listPattern.ReplaceAllStringFunc(input, func(m string) string {
		content := listPattern.FindStringSubmatch(m)[1]
		items := splitListItems(content)
		var buf strings.Builder
		buf.WriteString("<ul>")
		for _, it := range items {
			buf.WriteString("<li>" + parseLists(it) + "</li>")
		}
		buf.WriteString("</ul>")
		return buf.String()
	})

	// 处理有序列表
	input = olistPattern.ReplaceAllStringFunc(input, func(m string) string {
		content := olistPattern.FindStringSubmatch(m)[1]
		items := splitListItems(content)
		var buf strings.Builder
		buf.WriteString("<ol>")
		for _, it := range items {
			buf.WriteString("<li>" + parseLists(it) + "</li>")
		}
		buf.WriteString("</ol>")
		return buf.String()
	})

	return input
}

// 安全拆分列表项
func splitListItems(content string) []string {
	// 把 [*] 替换成一个特殊分隔符
	tmp := strings.ReplaceAll(content, "[*]", "\x00")
	parts := strings.Split(tmp, "\x00")
	var items []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			items = append(items, p)
		}
	}
	return items
}

// Markdown 转 HTML
func MarkdownToHTML(markdownContent string) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	err := md.Convert([]byte(markdownContent), &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
