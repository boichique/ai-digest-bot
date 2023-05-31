from http.server import BaseHTTPRequestHandler, HTTPServer
import json
import argparse
import os
import pytube
import whisper

class RequestHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path == '/transcribe':
            content_length = self.headers.get('Content-Length')
            if content_length is None:
                self.send_error(400, 'Content-Length header is missing')
                return 
            content_length = int(content_length)
            post_data = self.rfile.read(content_length)
            data = json.loads(post_data.decode('utf-8'))
            link = data['link']
            output = data['output']

            # Downloading the audio file from YouTube
            yt = pytube.YouTube(link)
            audio = yt.streams.filter(only_audio=True).first()
            audio.download(filename=output+'.mp4')
            
            # Transcribing the audio file
            model = whisper.load_model('base')
            text = model.transcribe(output+'.mp4')
            os.remove(output+'.mp4')

            # Sending the transcription as the response
            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()
            self.wfile.write(text['text'].encode('utf-8'))
        else:
            self.send_error(404)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Transcriptor youtube video to text')
    parser.add_argument('--port', type=int, default=10001, help='Port to listen on')
    args = parser.parse_args()
    server_address = ('', args.port)
    httpd = HTTPServer(server_address, RequestHandler)
    print('Listening on port', args.port)
    httpd.serve_forever()
