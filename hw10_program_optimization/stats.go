package hw10programoptimization

import (
	"bytes"
	"context"
	"errors"
	"io"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

//easyjson:json
type User struct {
	ID       int
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

const (
	readSize        int = 1024 * 4
	stringsChanSize int = 300
)

var (
	nDomainsWorkers = runtime.NumCPU()
	errNilReader    = errors.New("io.reader is nil")
)

var userPool = sync.Pool{
	New: func() interface{} { return new(User) },
}

type syncDomains struct {
	sync.Mutex
	domains DomainStat
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, errNilReader
	}

	wg, ctx := errgroup.WithContext(context.Background())

	byteCh := readStep(ctx, r, wg)
	strCh := splitStep(ctx, byteCh, wg)
	res := getDomainsStep(ctx, domain, strCh, wg)

	err := wg.Wait()

	return res, err
}

func readStep(ctx context.Context, r io.Reader, wg *errgroup.Group) <-chan []byte {
	resCh := make(chan []byte, 1)
	wg.Go(func() error {
		defer close(resCh)

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				readEnd := false
				buf := make([]byte, readSize)
				n, err := io.ReadAtLeast(r, buf, readSize)
				if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
					readEnd = true
				} else if err != nil {
					return err
				}

				resCh <- buf[:n]

				if readEnd {
					return nil
				}
			}
		}
	})

	return resCh
}

func splitStep(ctx context.Context, byteCh <-chan []byte, wg *errgroup.Group) <-chan []byte {
	resCh := make(chan []byte, stringsChanSize)

	wg.Go(func() error {
		defer close(resCh)

		residue := make([]byte, 0)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case page, ok := <-byteCh:
				if !ok {
					if len(residue) != 0 {
						resCh <- residue
					}

					return nil
				}

				lines := bytes.Split(page, []byte{'\n'})

				lines[0] = append(residue, lines[0]...)
				residue = writeLinesToChanel(lines, resCh)
			}
		}
	})

	return resCh
}

func writeLinesToChanel(lines [][]byte, ch chan<- []byte) []byte {
	for i, l := range lines {
		if i == len(lines)-1 {
			if len(l) > 0 && l[len(l)-1] != '\n' {
				return l
			}
		}

		ch <- l
	}

	return make([]byte, 0)
}

func getDomainsStep(ctx context.Context, domain string, byteCh <-chan []byte, wg *errgroup.Group) DomainStat {
	res := syncDomains{sync.Mutex{}, DomainStat{}}
	rg, err := regexp.Compile("\\." + domain)
	if err != nil {
		wg.Go(func() error {
			return err
		})
	}

	for i := 0; i < nDomainsWorkers; i++ {
		wg.Go(func() error {
			return getDomainWorker(ctx, byteCh, rg, &res)
		})
	}

	return res.domains
}

func getDomainWorker(ctx context.Context, byteCh <-chan []byte, rg *regexp.Regexp, domains *syncDomains) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case str, ok := <-byteCh:
			if !ok {
				return nil
			}

			if err := getDomain(rg, str, domains); err != nil {
				return err
			}
		}
	}
}

func getDomain(rg *regexp.Regexp, rowBytes []byte, domains *syncDomains) error {
	u := userPool.Get().(*User)
	defer userPool.Put(u)
	u.Email = ""

	err := getUser(rowBytes, u)
	if err != nil {
		return err
	}

	checkDomain(rg, u, domains)

	return nil
}

func getUser(rowBytes []byte, u *User) error {
	if err := u.UnmarshalJSON(rowBytes); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func checkDomain(rg *regexp.Regexp, u *User, domains *syncDomains) {
	if matched := rg.Match([]byte(u.Email)); matched {
		domains.Lock()
		num := domains.domains[strings.ToLower(strings.SplitN(u.Email, "@", 2)[1])]
		num++
		domains.domains[strings.ToLower(strings.SplitN(u.Email, "@", 2)[1])] = num
		domains.Unlock()
	}
}
