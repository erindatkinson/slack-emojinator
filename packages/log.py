"""module for logging configuration"""
import structlog


def get_logger():
    """get the logger"""

    #TODO: add config
    return structlog.get_logger()
