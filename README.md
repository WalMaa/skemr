<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->

<a id="readme-top"></a>

<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![project_license][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">

<h3 align="center">Skemr</h3>

  <p align="center">
    Validate custom schema rules in CI/CD pipelines
    <br />
    <a href="https://github.com/walmaa/skemr/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    &middot;
    <a href="https://github.com/walmaa/skemr/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>

<!-- ABOUT THE PROJECT -->

## About The Project

![Skemr database visualization][product-screenshot]

DOCS UNDER CONSTRUCTION

<!-- GETTING STARTED -->

## Quick Start

### Docker

```bash
docker compose -f docker-compose.yml up -d --build
```

<!-- USAGE EXAMPLES -->

## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/walmaa/skemr.svg?style=for-the-badge
[contributors-url]: https://github.com/walmaa/skemr/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/walmaa/skemr.svg?style=for-the-badge
[forks-url]: https://github.com/walmaa/skemr/network/members
[stars-shield]: https://img.shields.io/github/stars/walmaa/skemr.svg?style=for-the-badge
[stars-url]: https://github.com/walmaa/skemr/stargazers
[issues-shield]: https://img.shields.io/github/issues/walmaa/skemr.svg?style=for-the-badge
[issues-url]: https://github.com/walmaa/skemr/issues
[license-shield]: https://img.shields.io/github/license/walmaa/skemr.svg?style=for-the-badge
[license-url]: https://github.com/walmaa/skemr/blob/master/LICENSE.txt
[product-screenshot]: images/screenshot.png

## CLI

Run validate command with docker cli

```bash
 docker run walmaa/skemr-cli validate -P aaaabbbb-aaaa-aaaa-aaaa-aaaabbbbcccc -D 11112222-3333-4444-5555-666677778888 -T SuHFpgA8YOon.CF2heO2_0_6Xjf2ilfWPiOADTNq2y25Lqg-CeBHRpVk --host http://host.docker.internal:8080 -M .
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>
