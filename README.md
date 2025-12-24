# Arch Linux News (alnews)

A simple script in Go to fetch and display the latest news from the Arch Linux website.

## Features

- Fetches the latest news from the Arch Linux website.
- Displays news titles and links in the terminal.
- Fuzzy search functionality to filter news by keywords.

## Installation

You can install `alnews` by cloning the repository and building it with Go:

```bash
git clone https://github.com/hsm-gustavo/alnews.git
cd alnews
go build -o alnews
```

Alternatively, you can download pre-built binaries from the [releases](https://github.com/hsm-gustavo/alnews/releases) page.

## Requirements

- Go 1.24.0 or higher
- Internet connection to fetch news
- Optional:
  - xdg-open to open links in the default browser.

## Usage

Run the `alnews` executable from your terminal:

```bash
./alnews [options]
```

### Example

```bash
./alnews -h

Output:

A simple script written in Go to fetch and display the latest news from the Arch Linux website, with colorful output, date and link for the post.

Usage:
  alnews [flags]

Flags:
  -h, --help            help for alnews
  -i, --inspect int8    a longer print of the news (index, starting from 0) (default -1)
  -l, --limit uint8     number of news to show (default 3)
  -o, --open int8       opens the link of a specified news (index, starting from 0) in your default browser (default -1)
  -r, --refresh         if you want to refresh data or use cache
  -s, --search string   uses fuzzy search to filter news titles for specific keywords (e.g., 'nvidia')
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
