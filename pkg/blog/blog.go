package blog

import (
	"errors"
	"pokapoka-viewer/pkg/utils"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
 * Blog struct
 */
type Blog struct {
	Title string
	Detail string
	MediaURLs []string
	Comments []string
}

/*
 * This function is used to get the blog.
 * called from input handler
 */
func GetBlog(url string) (*Blog, error) {

	blogType, err := utils.ParseURL(url)

	if err != nil {
		return nil, err
	}

	switch blogType {
	case utils.Poka:
		return getPokaBlog(url)
	case utils.Carro:
		return getCarroBlog(url)
	case utils.Grotty:
		return getGrottyBlog(url)
	default:
		return nil, nil
	}
}


func getPokaBlog(url string) (*Blog, error) {
	htmlText, err := utils.FetchHTML(url)
	if err != nil {
		return nil, err
	}

	titlePattern := `<meta name="twitter:title" content="(.*?)" />`
	detailPattern := `<meta name="twitter:description" content="(.*?)" />`

	reTitle := regexp.MustCompile(titlePattern)
	reDetail := regexp.MustCompile(detailPattern)

	matchesTitle := reTitle.FindStringSubmatch(htmlText)
	matchesDetail := reDetail.FindStringSubmatch(htmlText)

	if len(matchesTitle) < 2 {
		return nil, errors.New("no title found")
	}

	if len(matchesDetail) < 2 {
		return nil, errors.New("no detail found")
	}

	title := matchesTitle[1]
	detail := matchesDetail[1]


	var mediaUrls []string
	mediaPattern1 := `<video class="wp-video-shortcode".*?動画が見れない方はこちら</a>`
    mediaPattern2 := `href="(https://.*?)"`

	reMediaPattern1 := regexp.MustCompile(mediaPattern1)
	reMediaPattern2 := regexp.MustCompile(mediaPattern2)

	matchesMedia1 := reMediaPattern1.FindAllString(htmlText, -1)
	for _, match := range matchesMedia1 {
		matchesMedia2 := reMediaPattern2.FindAllStringSubmatch(match, -1)
		mediaUrls = append(mediaUrls, matchesMedia2[1][1])
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))

	if err != nil {
		return nil, err
	}

	var comments []string
	doc.Find("div.comments-area li.comment-set").Each(func(i int, s *goquery.Selection) {
		text := s.Find("p").Eq(1).Text()
		comments = append(comments, text)
	})

	return &Blog{
		Title: title,
		Detail: detail,
		MediaURLs: mediaUrls,
		Comments: comments,
	}, nil
}

func getCarroBlog(url string) (*Blog, error) {
	htmlText, err := utils.FetchHTML(url)
	if err != nil {
		return nil, err
	}

	titlePattern := `<meta name="twitter:title" content="(.*?)" />`
	detailPattern := `<meta name="twitter:description" content="(.*?)" />`

	reTitle := regexp.MustCompile(titlePattern)
	reDetail := regexp.MustCompile(detailPattern)

	matchesTitle := reTitle.FindStringSubmatch(htmlText)
	matchesDetail := reDetail.FindStringSubmatch(htmlText)

	if len(matchesTitle) < 2 {
		return nil, errors.New("no title found")
	}

	if len(matchesDetail) < 2 {
		return nil, errors.New("no detail found")
	}

	title := matchesTitle[1]
	detail := matchesDetail[1]


	mediaPattern1 := `<a href=".*?</a>\s*<p class="playback_movie">└画像クリックで別タブで動画を再生します┘</p>`
	mediaPattern2 := `href="(https://.*?)"` 

	var mediaUrls []string

	reMediaPattern1 := regexp.MustCompile(mediaPattern1)
	reMediaPattern2 := regexp.MustCompile(mediaPattern2)

	matchesMedia1 := reMediaPattern1.FindAllString(htmlText, -1)
	for _, match := range matchesMedia1 {
		matchesMedia2 := reMediaPattern2.FindAllStringSubmatch(match, -1)
		mediaUrls = append(mediaUrls, matchesMedia2[0][1])
	}

	var comments []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))

	if err != nil {
		return nil, err
	}

	doc.Find("ol.commentlist li.comment").Each(func(i int, s *goquery.Selection) {
		text := s.Find("p").First().Text()
		comments = append(comments, text)
	})

	return &Blog{
		Title: title,
		Detail: detail,
		MediaURLs: mediaUrls,
		Comments: comments,
	}, nil
}

func getGrottyBlog(url string) (*Blog, error) {
	htmlText, err := utils.FetchHTML(url)
	if err != nil {
		return nil, err
	}

	titlePattern := `<meta name="twitter:title" content="(.*?)" />`
	detailPattern := `<meta name="twitter:description" content="(.*?)" />`

	reTitle := regexp.MustCompile(titlePattern)
	reDetail := regexp.MustCompile(detailPattern)

	matchesTitle := reTitle.FindStringSubmatch(htmlText)
	matchesDetail := reDetail.FindStringSubmatch(htmlText)

	if len(matchesTitle) < 2 {
		return nil, errors.New("no title found")
	}

	if len(matchesDetail) < 2 {
		return nil, errors.New("no detail found")
	}

	title := matchesTitle[1]
	detail := matchesDetail[1]

	pattern1 := `<a href=".*?" target="_blank" rel="noreferrer nofollow">動画に行くｗ</a></p>`
	pattern2 := `href="(https://.*?)"`

	var mediaUrls []string
	matchesMedia1 := regexp.MustCompile(pattern1).FindAllString(htmlText, -1)
	for _, match := range matchesMedia1 {
		matchesMedia2 := regexp.MustCompile(pattern2).FindAllStringSubmatch(match, -1)
		if len(matchesMedia2) < 1 {
			return nil, errors.New("no media found")
		}
		mediaUrls = append(mediaUrls, matchesMedia2[0][1])
	}

	var comments []string
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlText))

	if err != nil {
		return nil, err
	}

	doc.Find("#comment_list li.comment").Each(func(i int, s *goquery.Selection) {
		text := s.Find("p").First().Text()
		comments = append(comments, text)
	})

	return &Blog{
		Title: title,
		Detail: detail,
		MediaURLs: mediaUrls,
		Comments: comments, 
	}, nil
}