# WebRTC Signaling Server Ayame

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/OpenAyame/ayame.svg)](https://github.com/OpenAyame/ayame)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Actions Status](https://github.com/OpenAyame/ayame/workflows/Go%20Build%20&%20Format/badge.svg)](https://github.com/OpenAyame/ayame/actions)

## About Shiguredo's open source software

We will not respond to PRs or issues that have not been discussed on Discord. Also, Discord is only available in Japanese.

Please read https://github.com/shiguredo/oss/blob/master/README.en.md before use.

## 時雨堂のオープンソースソフトウェアについて

利用前に https://github.com/shiguredo/oss をお読みください。

## WebRTC Signaling Server Ayame について

WebRTC Signaling Server Ayame は WebRTC 向けのシグナリングサーバです。

WebRTC の P2P でのみ動作します。また動作を 1 ルームを最大 2 名に制限することでコードを小さく保っています。

## OpenAyame プロジェクトについて

OpenAyame は WebRTC Signaling Server Ayame をオープンソースとして公開し、
継続的に開発を行うことで WebRTC をより身近に、使いやすくするプロジェクトです。

詳細については下記をご確認ください。

[OpenAyame プロジェクト](http://bit.ly/OpenAyame)

## 方針

- シグナリングの仕様の破壊的変更を可能な限り行わない
- Go のバージョンは定期的にアップデートを行う
- 依存ライブラリは定期的にアップデートを行う

## 注意

- Ayame は P2P にしか対応していません
- Ayame は 1 ルーム最大 2 名までしか対応していません

## 使ってみる

Ayame を使ってみたい人は [USE.md](docs/USE.md) をお読みください。

## Web SDK を使ってみる

[Ayame Web SDK](https://github.com/OpenAyame/ayame-web-sdk)

## Web SDK サンプルを使ってみる

[Ayame Web SDK サンプル](https://github.com/OpenAyame/ayame-web-sdk-samples)

## 仕組みの詳細を知りたい

Ayame の仕組みを知りたい人は [OpenAyame/ayame\-spec](https://github.com/OpenAyame/ayame-spec) をお読みください。

## Ayame Labo を使ってみる

Ayame 仕様と完全互換な STUN/TURN サーバやルーム認証を組み込んだ無料で利用可能なシグナリングサービスを時雨堂が提供しています。

[Ayame Labo](https://ayame-labo.shiguredo.app/)

## ライセンス

Apache License 2.0

```
Copyright 2019-2022, Shiguredo Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Ayame 利用例

- [tarakoKutibiru/UnityRenderStreaming\-Ayame\-Sample](https://github.com/tarakoKutibiru/UnityRenderStreaming-Ayame-Sample)
- [tarukosu/MixedReality\-WebRTC\-ayame: MixedReality\-WebRTC にて、シグナリングサーバとして Ayame を利用するためのコード](https://github.com/tarukosu/MixedReality-WebRTC-ayame)
- [MixedReality\-WebRTC と Ayame Labo を利用して Unity で WebRTC を使う](https://zenn.dev/tarukosu/articles/20210220-webrtc-ayame)
- [kadoshita/kisei\-online: 手軽に使える，オンライン帰省用ビデオ通話ツール](https://github.com/kadoshita/kisei-online)
- [hakobera/serverless\-webrtc\-signaling\-server: Serverless WebRTC Signaling Server only works for WebRTC P2P\.](https://github.com/hakobera/serverless-webrtc-signaling-server)
- [mganeko/react\_ts\_ayame: React\.js and Typescript example for Ayame Labo \(WebRTC signaling\)](https://github.com/mganeko/react_ts_ayame)
- [mganeko/react\_ts\_ayame\_recv: React\.js and Typescript example for Ayame Labo \(WebRTC signaling\)](https://github.com/mganeko/react_ts_ayame_recv)
