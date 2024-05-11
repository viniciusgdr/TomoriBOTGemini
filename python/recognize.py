import asyncio
import sys
from shazamio import Shazam
import json


async def main():
  shazam = Shazam()
  out = await shazam.recognize_song(sys.argv[1])
  r = json.dumps(out)
  print(r)

loop = asyncio.get_event_loop()
loop.run_until_complete(main())