import argparse
import os
import pytube
import whisper

# Creating an argument parser to accept the video link and audiofile name
parser = argparse.ArgumentParser(description='Transcribe audio from a YouTube video')
parser.add_argument('link', metavar='link', type=str, help='YouTube video link')
parser.add_argument('output', metavar='output', type=str, default='audio.mp4', help='Output file name')
args = parser.parse_args()

# Reading the above Taken movie Youtube link
data = pytube.YouTube(args.link)

# Converting and downloading as 'MP4' file
audio = data.streams.get_audio_only()
audio.download(filename=args.output+'.mp4')

model = whisper.load_model('base')
text = model.transcribe(args.output+'.mp4')

# Printing the transcribe
print(text['text'])

os.remove(args.output+'.mp4')
