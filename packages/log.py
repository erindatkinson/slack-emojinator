"""module for logging configuration"""

import structlog


# TODO: Add configuration plumbing
# def init_logger():
#     """Configure"""
#     structlog.configure_once()


def get_logger():
    """get the logger"""
    return structlog.get_logger()
