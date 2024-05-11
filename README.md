# TomoriBOT WhatsApp - Alimentado por IA
TomoriBOT é um bot de WhatsApp alimentado por IA, que utiliza a API da Gemini (Google) para responder os usuários de forma inteligente e natural. Nele, você pode interagir com o bot através de mensagens de texto e baixar músicas do YouTube de forma rápida e prática.
<img src="/public/image-banner.png" alt="Banner"/>

## Soluções/Utilidades
Sabemos que a maioria das pessoas utilizam o WhatsApp para se comunicar com amigos e familiares, mas também para se entreter. Pensando nisso, o TomoriBOT foi criado para facilitar a vida dos usuários, trazendo uma experiência única e inovadora. Com ele, você pode baixar músicas do YouTube, baixar vídeos do Twitter, Instagram e TikTok, criar figurinhas, jogar uma moeda e muito mais.

Além disso, muitas tarefas que antes eram feitas manualmente, agora podem ser feitas de forma automática, economizando tempo e esforço. Como por exemplo, baixar músicas do YouTube, que antes era necessário acessar um site, colar o link, esperar o download e por fim, baixar a música. Com o TomoriBOT, você pode fazer isso de forma rápida e prática, apenas enviando o nome da música que deseja baixar.

## Recursos
- ✅ Buscar músicas do YouTube
- ✅ Baixar em MP3/MP4 qualquer conteúdo do YouTube
- ✅ Baixar vídeos do Twitter.
- ✅ Conversar com o bot
- ✅ Criar figurinhas
- ✅ Jogue uma Moeda com o bot. (Cara ou Coroa)
- ✅ Reconhecer músicas (Shazam)
- ✅ Baixar vídeos do Instagram (Reels)
- ✅ Baixar vídeos do TikTok

## Eficiência	
O TomoriBOT foi desenvolvido para ser eficiente e rápido, permitindo que os usuários interajam com o bot de forma rápida e prática. Com ele, você pode baixar músicas do YouTube em poucos segundos, sem precisar acessar sites ou instalar programas. Além disso, o bot é capaz de reconhecer músicas em tempo real, permitindo que você descubra o nome da música que está tocando no ambiente.

- 🚀 Suporta Chats Privados
- 🚀 Suporta Grupos (Digitando "Tomori," na frente)

## Instalação
Para instalar o TomoriBOT, você precisa ter o Node.js, Python, FFmpeg e o Go instalados em sua máquina. Após isso, basta clonar o repositório e instalar as dependências.

```bash
# Clone o repositório
$ git clone https://github.com/viniciusgdr/TomoriBOTGemini
$ cd TomoriBOTGemini
$ bash installer.sh
```

## Configuração
Para configurar o TomoriBOT, você precisa criar um arquivo `.env` na raiz do projeto e adicionar as seguintes variáveis de ambiente.

```env
GEMINI_APIKEY=YOUR_API_KEY
PHONE_NUMBER=YOUR_PHONE_NUMBER
```
## Demonstração
<div style="display: flex; flex-wrap: wrap;gap: 3px">
  <img src="/public/img1.jpeg" alt="Demo" width="48%" />
  <img src="/public/img5.jpeg" alt="Demo" width="48%" />
  <img src="/public/img4.jpeg" alt="Demo" width="48%" />
  <img src="/public/img3.jpeg" alt="Demo" width="48%" />
  <img src="/public/img2.jpeg" alt="Demo" width="48%" />
</div>

## Arquitetura Robusta
O TomoriBOT foi desenvolvido com uma arquitetura robusta, utilizando a API da Gemini (Google) para responder os usuários de forma inteligente e natural. Além disso, ele foi desenvolvido usando os princípios do Clean Architecture, que permite a fácil manutenção e escalabilidade do projeto.
Os módulos em Python e Node foram feitos para acelerar o desenvolvimento de diversos recursos que a linguagem oferece, como reconhecimento de músicas (Em Python) e Manipulação de arquivos Webp (Em Node).
*A Lib Whatsmeow não foi desacoplada 100% pelo fato que não tem outra lib que faça o mesmo trabalho.

## Privacidade
O TomoriBOT respeita a privacidade dos usuários e não armazena nenhuma informação pessoal. Todas as mensagens trocadas com o bot são processadas em tempo real e não são armazenadas em nenhum banco de dados (ainda). Além disso, o bot não compartilha nenhuma informação com terceiros e não exibe anúncios. Todo conteúdo que é efetuado download de terceiros é provido de APIs públicas.