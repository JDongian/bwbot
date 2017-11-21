package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "io/ioutil"
    "net/url"
    "syscall"
    "strings"
    "github.com/bwmarrin/discordgo"
    "github.com/PuerkitoBio/goquery"
)

// Variables used for command line parameters
var (
    Token string
)

func init() {
    flag.Parse()
    var data []byte
    switch flag.NArg() {
    case 0:
        data, _ = ioutil.ReadAll(os.Stdin)
        break
    case 1:
        data, _ = ioutil.ReadFile(flag.Arg(0))
        break
    default:
        fmt.Printf("Supply token via stdin or file.\n")
        os.Exit(1)
    }
    Token = strings.TrimSpace(string(data))
}

func main() {
    fmt.Println("Starting...")
    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + Token)
    if err != nil {
        fmt.Println("error creating Discord session,", err)
        return
    }

    // Register the messageCreate func as a callback for MessageCreate events.
    dg.AddHandler(messageCreate)

    // Open a websocket connection to Discord and begin listening.
    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    // Wait here until CTRL-C or other term signal is received.
    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Cleanly close down the Discord session.
    dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Ignore all messages created by the bot itself
    if m.Author.ID == s.State.User.ID {
        return
    }
    msg := m.ContentWithMentionsReplaced()
    parts := strings.Split(strings.ToLower(msg), " ")
    // If the message starts with the prefix
    if parts[0] == "!bw" {
        handleMessage(s, m, parts[1:])
    }
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, tokens []string) {
    mapData := tlpdSearchResults(strings.Join(tokens, " "))
    em := formatMapData(mapData)
    s.ChannelMessageSendEmbed(m.ChannelID, em)
}

func tlpdSearchResults(query string) map[string]string {
    BaseUrl := "http://www.teamliquid.net"
    searchUrl := fmt.Sprintf(
        BaseUrl +
        "/tlpd/maps/index.php?" +
        "section=korean&tabulator_page=1&tabulator_order_col=default&" +
        "tabulator_search=%s",
        query)
    fmt.Println(searchUrl)
    doc, err := goquery.NewDocument(searchUrl)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("<HTML>" + doc.Text()[:20] + "...</HTML>")  // DEBUG

    // Find the map links.
    //doc.Find("div.roundcont #tblt_table table td a").Each(func(i int, s *goquery.Selection) {
    //    link, _ := s.Attr("href")
    //    fmt.Printf(link + "\n")
    //})
    first_map_link, _ := doc.Find("div.roundcont #tblt_table tr td a").First().Attr("href")
    fmt.Println(BaseUrl + first_map_link)  // Debug only
    return tlpdParseMapLink(BaseUrl + first_map_link)
}

func tlpdParseMapLink(link string) map[string]string {
    m := make(map[string]string)
    doc, err := goquery.NewDocument(link)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("<HTML>" + doc.Text()[:20] + "...</HTML>")  // DEBUG
    winRates := tlpdGetMapWinRates(doc)
    for k, v := range winRates {
        m[k] = v
    }
    imageLink := tlpdGetMapImageLink(doc)
    m["imageLink"] = imageLink
    m["link"] = link
    return m
}

func tlpdGetMapWinRates(doc *goquery.Document) map[string]string {
    data := doc.Find("div.roundcont table tbody td")
    return map[string]string{
        "tvzGames": data.Eq(1).Text(),
        "tvzWR": data.Eq(2).Text(),
        "zvpGames": data.Eq(5).Text(),
        "zvpWR": data.Eq(6).Text(),
        "pvtGames": data.Eq(9).Text(),
        "pvtWR": data.Eq(10).Text(),
    }
}

func tlpdGetMapImageLink(doc *goquery.Document) string {
    BaseUrl := "http://www.teamliquid.net"
    imageLink, _ := doc.Find("div.roundcont p a img").First().Attr("src")
    return BaseUrl + strings.TrimSpace(imageLink)
}

func formatMapData(mapData map[string]string) *discordgo.MessageEmbed {
    winRates := fmt.Sprintf(
        "TvZ: %s %s\n" +
        "ZvP: %s %s\n" +
        "PvT: %s %s",
        mapData["tvzGames"],
        mapData["tvzWR"],
        mapData["zvpGames"],
        mapData["zvpWR"],
        mapData["pvtGames"],
        mapData["pvtWR"],)
    link := (&url.URL{Path: mapData["link"]}).String()[2:]
    imageLink := (&url.URL{Path: mapData["imageLink"]}).String()[2:]

    fmt.Println("link: " + link)  // DEBUG
    fmt.Println("winrates: " + winRates)  // DEBUG
    fmt.Println("img: " + imageLink)  // DEBUG
    em := &discordgo.MessageEmbed{
        URL: link,
        Title: "TLPD Map Database",
        Description: winRates,
        Color: 0x254673,
        Image: &discordgo.MessageEmbedImage{
            URL: imageLink,
        },
    }
    return em
}
