import grpc

from .proto import sentence_gen_pb2, sentence_gen_pb2_grpc

SERVER_ADDR = "sentence-gen-grpc-server-218429278955.europe-north1.run.app:443"
TIMEOUT = 50  # seconds, matches Go's timeOut = time.Second * 50


# Error types matching Go rpc/client.go
class GrpcError(Exception):
    pass

class InvalidArgumentError(GrpcError):
    pass

class InternalServerError(GrpcError):
    pass

class ResourceExhaustedError(GrpcError):
    pass

class DeadlineExceededError(GrpcError):
    pass

class UnavailableError(GrpcError):
    pass

class UnknownError(GrpcError):
    pass


class GrpcClient:
    def __init__(self):
        credentials = grpc.ssl_channel_credentials()
        self._channel = grpc.secure_channel(SERVER_ADDR, credentials)
        self._stub = sentence_gen_pb2_grpc.SentenceGenStub(self._channel)

    def close(self):
        self._channel.close()

    def generate_sentence(self, word, word_language, translation_language,
                          translation_hint, include_audio, voice_gender):
        try:
            req = sentence_gen_pb2.GenerateSentenceRequest(
                word=word,
                word_language=word_language,
                translation_language=translation_language,
                translation_hint=translation_hint,
                include_audio=include_audio,
                voice_gender=voice_gender,
            )
            return self._stub.GenerateSentence(req, timeout=TIMEOUT)
        except grpc.RpcError as e:
            raise self._map_error(e) from e

    def translate(self, word, from_language, to_language,
                  translation_hint, include_audio, voice_gender):
        try:
            req = sentence_gen_pb2.TranslateRequest(
                word=word,
                from_language=from_language,
                to_language=to_language,
                translation_hint=translation_hint,
                include_audio=include_audio,
                voice_gender=voice_gender,
            )
            return self._stub.Translate(req, timeout=TIMEOUT)
        except grpc.RpcError as e:
            raise self._map_error(e) from e

    def generate_definition(self, word, language, definition_hint,
                            include_audio, voice_gender):
        try:
            req = sentence_gen_pb2.GenerateDefinitionRequest(
                word=word,
                language=language,
                definition_hint=definition_hint,
                include_audio=include_audio,
                voice_gender=voice_gender,
            )
            return self._stub.GenerateDefinition(req, timeout=TIMEOUT)
        except grpc.RpcError as e:
            raise self._map_error(e) from e

    def _map_error(self, e):
        code = e.code()
        msg = e.details() or ""
        if code == grpc.StatusCode.INVALID_ARGUMENT:
            return InvalidArgumentError(msg)
        if code == grpc.StatusCode.DEADLINE_EXCEEDED:
            return DeadlineExceededError(msg)
        if code == grpc.StatusCode.RESOURCE_EXHAUSTED:
            return ResourceExhaustedError(msg)
        if code == grpc.StatusCode.INTERNAL:
            return InternalServerError(msg)
        if code == grpc.StatusCode.UNAVAILABLE:
            return UnavailableError(msg)
        return UnknownError(msg)
