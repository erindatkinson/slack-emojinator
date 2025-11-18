"""module for logging configuration"""

import logging
import structlog



def get_logger():
    """get the logger"""
    structlog.configure(
        wrapper_class=structlog.make_filtering_bound_logger(logging.INFO),
    )
    return structlog.get_logger()
