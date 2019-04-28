from concurrent import futures
from importlib import import_module
import argparse
import time

import grpc

import gamer_pb2
import gamer_pb2_grpc


_DEFAULT_TIMEOUT = 60


class GamerServer(gamer_pb2_grpc.GamerServicer):
    def __init__(self, action):
        self.action = action
        self.userdata = {}

    def Ping(self, request, context):
        return gamer_pb2.PingMessage(id='dummy')

    def Action(self, request, context):
        context = {
            'gameId': request.gameId,
            'totalRound': request.totalRound,
            'currentRound': request.currentRound,
            'userdata': self.userdata,
        }
        result, userdata = self.action(request.data, context)

        self.userdata = userdata
        return gamer_pb2.ActionOutput(
            id=request.id,
            data=result,
        )


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--port')
    parser.add_argument('--src')

    args = parser.parse_args()

    # Load action codes
    mod = import_module(args.src)
    action = getattr(mod, 'action')

    # Start server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=2))
    gamer_pb2_grpc.add_GamerServicer_to_server(GamerServer(action), server)

    server.add_insecure_port(f"127.0.0.1:{args.port}")
    server.start()

    time.sleep(_DEFAULT_TIMEOUT)
    server.stop(0)


if __name__ == '__main__':
    main()
