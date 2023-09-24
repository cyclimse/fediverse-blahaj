INSERT INTO crawl_errors (error_code, description)
VALUES ('unknown', 'Unknown error'),
    ('timeout', 'Crawl timed out'),
    (
        'domain_not_found',
        'Domain could not be resolved'
    ),
    ('unreachable', 'Domain could not be reached'),
    (
        'invalid_nodeinfo',
        'Domain returned an invalid nodeinfo response'
    ),
    (
        'nodeinfo_version_not_supported_by_crawler',
        'Domain returned an unsupported nodeinfo version'
    ),
    (
        'invalid_json',
        'Domain returned an invalid json response'
    ),
    (
        'blocked_by_robots_txt',
        'Crawler was blocked by robots.txt'
    ),
    (
        'software_not_supported_by_crawler',
        'Software name is not supported by the crawler'
    ),
    (
        'software_version_not_supported_by_crawler',
        'Software version is not supported by the crawler'
    ),
    (
        'internal_error',
        'Internal error in the crawler'
    );