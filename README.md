# FireDorks

<div align="center">
  <img src="images/firedorks.png" width="300px">
</div>

***`fireDorks` is a small tool for scraping from Google Search results which can be used with AWS API Gateway.***

## Introduction

The tool performs queries on Google search engine to extract specific information, such as links, articles, or specific files using pattern parameter making it useful for OSINT (Open Source Intelligence) research.
To avoid being blocked by Google, the tool is designed to be used with [Fireprox](https://github.com/ustayready/fireprox). FireProx leverages the AWS API Gateway to create pass-through proxies that rotate the source IP address with every request.

### Benefits of using FireProx
 * Rotates IP address with every request
 * Configure separate regions
 * All HTTP methods supported
 * All parameters and URI's are passed through
 * Create, delete, list, or update proxies
 * Spoof X-Forwarded-For source IP header by requesting with an X-My-X-Forwarded-For header

## Usage
```shell
Usage: FireDorks --proxyurl="https://www.google.com" --dork=STRING

A tool for scraping Google Search results by @froyo75

Flags:
  -h, --help                                 Show context-sensitive help.
  -x, --proxyurl="https://www.google.com"    The proxy URL to use (default: https://www.google.com).
  -k, --dork=STRING                          The dork query to provide. (example: site:microsoft.com).
  -p, --maxpages=1000                        Specify the maximum number of Google search pages (default:1000).
  -r, --maxresults=100                       Specify the number of results per page to return (default:100).
  -w, --maxworkers=5                         Specify the number of workers to process queries (default:5).
  -s, --regex=STRING                         Specify a pattern to search using the POSIX standard regex.
  -e, --rawonly                              Specify whether values should be extracted from raw HTTP response data using a provided pattern.
  -l, --linksonly                            Specify whether to return only links in the output results (default:false).
  -f, --format="txt"                         Specify the output format (txt/csv/json) (default:txt).
  -d, --delay=2                              Specify the delay (in seconds) to avoid being fingerprinting (default:2).
  -c, --outconsole                           Specify whether to output all results to the console (default:false).
  -o, --outpath=STRING                       Output results to a specific file.
  -v, --version                              Print version information and quit.
```

**Create a new API gateway proxy in a specific AWS region (example: us-west-1) using https://www.google.com URL end-point**

```shell
python3 fire.py --access_key <ACCESS_KEY> --secret_access_key <SECRET_ACCESS_KEY> --region us-west-1 --command create --url https://www.google.com
```

**Retrieve all LinkedIn email addresses from a Google search using 10 threads with a maximum of 1000 pages for `mycompany.com`**
```shell
./firedorks -k "site:linkedin.com/in intext:@mycompany.com" -s "[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}" -c -p 1000 -r 100 -d 10 -e -o emails.list -x https://xxxxxxxxxx.execute-api.us-west-1.amazonaws.com/apps/

```

**List all documents with the `pdf` extension from Google search using 10 workers with a maximum of 1000 pages for `microsoft.com`**

```shell
./firedorks -k "filetype:pdf site:microsoft.com" -l -p 1000 -r 100 -d 10 -o results.txt -x https://xxxxxxxxxx.execute-api.us-west-1.amazonaws.com/apps/
```

**Retrieve all pages from Google search containing `*.passwords.txt` using 10 threads with a maximum of 1000 pages**

```shell
./firedorks -k "intitle:index of *.passwords.txt" -p 1000 -r 100 -d 10 -f csv -o results.csv -x https://xxxxxxxxxx.execute-api.us-west-1.amazonaws.com/apps/
```
