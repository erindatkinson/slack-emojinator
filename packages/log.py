"""module for logging configuration"""

import structlog


def get_logger():
    """get the logger"""
    return structlog.get_logger()
