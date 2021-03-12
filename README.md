<h1 align="center">
  <img src="https://github.com/Dreamacro/clash/raw/master/docs/logo.png" alt="Clash" width="200">
  <br>fast Clash fork for Dreamacro-clash; add memory cache<br>
</h1>

<h4 align="center">A rule-based tunnel in Go.</h4>

<p align="center">
  <a href="https://github.com/Dreamacro/clash/actions">
    <img src="https://img.shields.io/github/workflow/status/Dreamacro/clash/Go?style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/Dreamacro/clash">
    <img src="https://goreportcard.com/badge/github.com/Dreamacro/clash?style=flat-square">
  </a>
  <a href="https://github.com/Dreamacro/clash/releases">
    <img src="https://img.shields.io/github/release/Dreamacro/clash/all.svg?style=flat-square">
  </a>
</p>

## Alter patch

<p>1 user golang ccache replace clash build-in lru.cache;</p>
<p>2 add ccache cache rules not scan all rules by every time;</p>
<p>3 add altner rules API but not all.</p>
<p>4 add nutsdb to save dns result.</p>
