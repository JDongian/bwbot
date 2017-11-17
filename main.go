package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "io/ioutil"
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
    tlpdSearchResults(strings.Join(tokens, " "))
    s.ChannelMessageSend(m.ChannelID, "ok, " + strings.Join(tokens, " "))
}

func tlpdSearchResults(query string) {
    BaseUrl := "http://www.teamliquid.net"
    SearchUrl := fmt.Sprintf(
        BaseUrl +
        "/tlpd/maps/index.php?" +
        "section=korean&tabulator_page=1&tabulator_order_col=default&" +
        "tabulator_search=%s",
        query)
    fmt.Printf(SearchUrl)
    doc, err := goquery.NewDocument(SearchUrl)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf(doc.Text()[:100] + "\n\n")

    // Find the map links.
    doc.Find("div.roundcont #tblt_table table td a").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Attr("href")
        fmt.Printf(link + "\n")
    })
}

func parseResults(html []byte) {
    return
}
