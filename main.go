package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// CodeSearchResult estrutura a resposta da API de busca de código do GitHub.
type CodeSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []struct {
		HTMLURL     string `json:"html_url"`
		TextMatches []struct {
			Fragment string `json:"fragment"`
		} `json:"text_matches"`
	} `json:"items"`
}

func extractDomain(rawURL string) string {
	// Se a URL começar com //, adiciona "http:" para possibilitar o parse.
	if strings.HasPrefix(rawURL, "//") {
		rawURL = "http:" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	// Caso o host contenha a porta, separe-a.
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		// Se der erro, pode ser que não haja porta.
		host = u.Host
	}
	return host
}

func main() {
	// Flags de linha de comando:
	// -q: query simples para a API do GitHub.
	// -r: regex para filtrar os resultados.
	// -m: modo de extração: "urls" ou "domains". Quando definido, a extração será feita com uma regex interna e filtrada com -r.
	// -d: delay entre requisições.
	// -s: silent, apenas exibe os resultados extraídos sem a URL do arquivo, garantindo resultados únicos.
	apiQuery := flag.String("q", "", "Query de busca para a API do GitHub (ex: mercadolivre)")
	regexStr := flag.String("r", "", "Regex para filtrar os resultados localmente (ex: mercadolivre)")
	mode := flag.String("m", "", "Modo de extração: 'urls' ou 'domains' (opcional)")
	delay := flag.Int("d", 2, "Delay em segundos entre requisições para evitar bloqueio")
	silent := flag.Bool("s", false, "Silent: somente exibe os resultados extraídos (únicos), sem a URL do arquivo")
	flag.Parse()

	if *apiQuery == "" {
		log.Fatal("Você deve fornecer uma query para a API com o parâmetro -q")
	}
	if *regexStr == "" {
		log.Fatal("Você deve fornecer uma regex para filtrar os resultados com o parâmetro -r")
	}

	// Se não estivermos usando modo, compilamos a regex para filtrar os trechos.
	var re *regexp.Regexp
	var err error
	if *mode == "" {
		re, err = regexp.Compile(*regexStr)
		if err != nil {
			log.Fatalf("Erro ao compilar a regex: %v", err)
		}
	} else {
		// Se estiver usando modo ("urls" ou "domains"), compilamos a regex de filtro que será aplicada
		// sobre cada URL ou domínio extraído.
		re, err = regexp.Compile(*regexStr)
		if err != nil {
			log.Fatalf("Erro ao compilar a regex de filtro: %v", err)
		}
		if *mode != "urls" && *mode != "domains" {
			log.Fatal("O modo (-m) deve ser 'urls' ou 'domains'")
		}
	}

	// Regex interna para extração de URLs.
	urlRegex := regexp.MustCompile(`((https?:\/\/|\/\/)[^\s"'<>]+)`)

	// Obtém a chave do GitHub da variável de ambiente, se disponível.
	githubKey := os.Getenv("GITHUB_KEY")

	perPage := 100 // Máximo permitido pela API.
	page := 1

	// Mapa para garantir resultados únicos quando o modo silent estiver ativado.
	uniqueResults := make(map[string]bool)

	// Loop de paginação.
	for {
		baseURL := "https://api.github.com/search/code"
		// A query deve ser simples para a API.
		q := url.QueryEscape(*apiQuery)
		apiURL := fmt.Sprintf("%s?q=%s&page=%d&per_page=%d", baseURL, q, page, perPage)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatalf("Erro ao criar requisição: %v", err)
		}
		req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
		if githubKey != "" {
			req.Header.Set("Authorization", "token "+githubKey)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("Erro na requisição: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Fatalf("Erro da API (status %d): %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatalf("Erro ao ler resposta: %v", err)
		}

		var result CodeSearchResult
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatalf("Erro ao decodificar JSON: %v", err)
		}

		// Se não houver itens, encerra a busca.
		if len(result.Items) == 0 {
			if !*silent {
				fmt.Println("Nenhum resultado encontrado ou fim dos resultados disponíveis.")
			}
			break
		}

		// Processa cada item retornado e aplica o filtro.
		for _, item := range result.Items {
			for _, tm := range item.TextMatches {
				if *mode == "" {
					// Sem modo, usa a regex passada para filtrar os trechos.
					matches := re.FindAllString(tm.Fragment, -1)
					for _, m := range matches {
						if *silent {
							// Se silent, exibe somente o resultado, garantindo que seja único.
							if !uniqueResults[m] {
								fmt.Println(m)
								uniqueResults[m] = true
							}
						} else {
							fmt.Printf("\033[34m%s\033[0m - \033[32m%s\033[0m\n", item.HTMLURL, m)
						}
					}
				} else {
					// Com modo, extrai URLs usando a regex interna.
					urls := urlRegex.FindAllString(tm.Fragment, -1)
					for _, u := range urls {
						if *mode == "domains" {
							domain := extractDomain(u)
							if domain != "" && re.MatchString(domain) {
								if *silent {
									if !uniqueResults[domain] {
										fmt.Println(domain)
										uniqueResults[domain] = true
									}
								} else {
									fmt.Printf("\033[34m%s\033[0m - \033[32m%s\033[0m\n", item.HTMLURL, domain)
								}
							}
						} else if *mode == "urls" {
							if re.MatchString(u) {
								if *silent {
									if !uniqueResults[u] {
										fmt.Println(u)
										uniqueResults[u] = true
									}
								} else {
									fmt.Printf("\033[34m%s\033[0m - \033[32m%s\033[0m\n", item.HTMLURL, u)
								}
							}
						}
					}
				}
			}
		}

		// A API do GitHub retorna no máximo 1000 resultados (10 páginas com 100 itens cada).
		if page*perPage >= result.TotalCount || page >= 10 {
			if !*silent {
				fmt.Println("Fim dos resultados disponíveis.")
			}
			break
		}

		page++
		time.Sleep(time.Duration(*delay) * time.Second)
	}
}
