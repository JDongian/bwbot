from IPython import embed
import requests
import re
from bs4 import BeautifulSoup
from urllib.parse import quote

BASE_URL = "http://www.teamliquid.net"
SEARCH_URL = BASE_URL + "/tlpd/maps/index.php?section=korean&tabulator_page=1&tabulator_order_col=default&tabulator_search={query}"

HEADERS = {
    'Host': 'www.teamliquid.net',
    'Connection': 'keep-alive',
    'Cache-Control': 'max-age=0',
    'Upgrade-Insecure-Requests': '1',
    'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8',
    'Accept-Encoding': 'gzip, deflate',
    'Accept-Language': 'en-US,en;q=0.8',
    'Cookie': '__cfduid=df74cf2ecffe2b487ad5d80d768f6212a1509307992; SID=6hjurn007m6pere60e7a6tt0v3',
}


def search(query):
    print(SEARCH_URL.format(query=query))
    return requests.get(SEARCH_URL.format(query=query))


def _parse_winrates(html):
    """
    <td><strong>TvZ</strong>:</td>
    <td>148-139</td>
    <td>(51.6%)</td>
    <td>[ <a href="/tlpd/korean/maps/237_Fighting_Spirit/games/TvZ">Games</a> ]</td>
    <td><strong>ZvP</strong>:</td>
    <td>134-117</td>
    <td>(53.4%)</td>
    <td>[ <a href="/tlpd/korean/maps/237_Fighting_Spirit/games/ZvP">Games</a> ]</td>
    <td><strong>PvT</strong>:</td>
    <td>143-137</td>
    <td>(51.1%)</td>
    <td>[ <a href="/tlpd/korean/maps/237_Fighting_Spirit/games/TvP">Games</a> ]</td>
    """
    soup = BeautifulSoup(html, 'html.parser')
    rows = [row.contents[0] for row in soup.select(".roundcont table td")]
    tvz_games, tvz_wr = rows[1], rows[2]
    zvp_games, zvp_wr = rows[5], rows[6]
    pvt_games, pvt_wr = rows[9], rows[10]
    return {'TvZ': "{} {}".format(tvz_games, tvz_wr),
            'ZvP': "{} {}".format(zvp_games, zvp_wr),
            'PvT': "{} {}".format(pvt_games, pvt_wr),
            'summary': "TvZ: {} {}\n"
                       "ZvP: {} {}\n"
                       "PvT: {} {}".format(tvz_games, tvz_wr,
                                           zvp_games, zvp_wr,
                                           pvt_games, pvt_wr)}


def _parse_image(html):
    soup = BeautifulSoup(html, 'html.parser')
    img_path = soup.select(".roundcont img")[0].attrs['src'].strip()
    return BASE_URL + quote(img_path)


def parse_map_link(link):
    """Return a dictionary containing extracted map information.
    """
    html = requests.get(link, headers=HEADERS).content
    return {'link': link,
            'image_link': _parse_image(html),
            'win_rates': _parse_winrates(html)}


def get_map_links(html):
    maps = re.findall("/tlpd/korean/maps/\d+_\w+", str(html))
    return [BASE_URL + m for m in maps]


def get_map_stats(query):
    html = search(query).content
    print(html[:100])
    first_map = get_map_links(html)[0]
    result = parse_map_link(first_map)
    print(result)
    return result


if __name__ == "__main__":
    result = get_map_stats("fighting spirit")
