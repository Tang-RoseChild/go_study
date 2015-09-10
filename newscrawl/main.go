package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"mahonia"

	"encoding/json"
	"flag"
	"net/http"
)

type TopNews struct {
	resp *http.Response
	Doc  *goquery.Document
	S    *goquery.Selection
}

type news struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type jsonNews struct {
	Kind string `json:"kind"`
	News []news `json:"news"`
}

func main() {
	addr := flag.String("http", ":20003", "host ip address like `hostip:port` ")
	flag.Parse()

	http.HandleFunc("/tops", getTops)
	http.HandleFunc("/tops/baidu", getBaiduTops)
	http.HandleFunc("/tops/wangyi", getWangyiTops)
	http.HandleFunc("/tops/fenghuang", getFenghuangTops)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		fmt.Println("server error : ", err.Error())
	}

}

func CreateDoc(char, url, selector string) (*TopNews, error) {
	if char == "utf-8" || char == "UTF-8" {
		doc, err := goquery.NewDocument(url)
		if err != nil {

			return nil, fmt.Errorf("when create from reader error %s ", err.Error())
		}
		s := doc.Find(selector)
		return &TopNews{
			resp: nil,
			Doc:  doc,
			S:    s,
		}, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if mahonia.GetCharset(char) == nil {

		return nil, fmt.Errorf("%s charset not suported \n", char)
	}

	dec := mahonia.NewDecoder(char)
	rd := dec.NewReader(resp.Body)

	doc, err := goquery.NewDocumentFromReader(rd)
	if err != nil {

		return nil, fmt.Errorf("when create from reader error %s ", err.Error())
	}
	s := doc.Find(selector)
	return &TopNews{
		resp: resp,
		Doc:  doc,
		S:    s,
	}, nil
}

func getTops(w http.ResponseWriter, r *http.Request) {
	jn := make([]jsonNews, 0)

	// baidu tops
	baidu, err := baiduTops()
	if err != nil {
		fmt.Println("baidu tops err : ", err.Error())
		return
	}
	jn = append(jn, jsonNews{"baidu", baidu})

	wangyi, err := wangyiTops()
	if err != nil {
		fmt.Println("baidu tops err : ", err.Error())
		return
	}
	jn = append(jn, jsonNews{"wangyi", wangyi})

	fenghuang, err := fenghuangTops()
	if err != nil {
		fmt.Println("fenghuang tops err : ", err.Error())
		return
	}
	jn = append(jn, jsonNews{"fenghuang", fenghuang})

	fmt.Println("jsonNews : ", jn)
	enc := json.NewEncoder(w)
	if err := enc.Encode(jn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getWangyiTops(w http.ResponseWriter, r *http.Request) {
	jn := make([]jsonNews, 0)

	wangyi, err := wangyiTops()
	if err != nil {
		fmt.Println("wangyi tops err : ", err.Error())
		return
	}
	jn = append(jn, jsonNews{"wangyi", wangyi})

	fmt.Println("jsonNews : ", jn)
	enc := json.NewEncoder(w)
	if err := enc.Encode(jn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getFenghuangTops(w http.ResponseWriter, r *http.Request) {
	jn := make([]jsonNews, 0)

	fenghuang, err := fenghuangTops()
	if err != nil {
		fmt.Println("fenghuang tops err : ", err.Error())
		return
	}
	jn = append(jn, jsonNews{"fenghuang", fenghuang})

	fmt.Println("jsonNews : ", jn)
	enc := json.NewEncoder(w)
	if err := enc.Encode(jn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getBaiduTops(w http.ResponseWriter, r *http.Request) {
	jn := make([]jsonNews, 0)

	// baidu tops
	baidu, err := baiduTops()
	if err != nil {
		fmt.Println("baidu tops err : ", err.Error())
		return
	}
	jn = append(jn, jsonNews{"baidu", baidu})

	fmt.Println("jsonNews : ", jn)
	enc := json.NewEncoder(w)
	if err := enc.Encode(jn); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func baiduTops() ([]news, error) {
	baidu, err := CreateDoc("gb2312", "http://news.baidu.com/", ".hotnews a")
	// m := getWangYiTops("gbk")
	if err != nil {
		fmt.Println(" baidu err is", err.Error())
		return nil, err
	}
	var baiduNews = make([]news, 0)
	baidu.S.Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			// fmt.Printf("title %s href %s \n", s.Text(), href)
			baiduNews = append(baiduNews, news{s.Text(), href})

		}
	})
	return baiduNews, nil
}

func wangyiTops() ([]news, error) {
	wangyi, err := CreateDoc("gbk", "http://news.163.com/", ".ns-wnews.mb20 a")
	if err != nil {
		fmt.Println(" wangyi err is", err.Error())
		return nil, err
	}
	var wangyiNews = make([]news, 0)
	wangyi.S.Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			// fmt.Printf("title %s href %s \n", s.Text(), href)
			wangyiNews = append(wangyiNews, news{s.Text(), href})
		}
	})
	return wangyiNews, nil
}

func fenghuangTops() ([]news, error) {
	fenghuang, err := CreateDoc("utf-8", "http://news.ifeng.com/", ".box_01 a")
	if err != nil {
		fmt.Println(" fenghuang err is", err.Error())
		return nil, err
	}
	var fenghuangNews = make([]news, 0)
	fenghuang.S.Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			// fmt.Printf("title %s href %s \n", s.Text(), href)
			fenghuangNews = append(fenghuangNews, news{s.Text(), href})
		}
	})
	return fenghuangNews, nil
}
