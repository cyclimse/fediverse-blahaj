/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Crawl = {
    id: string;
    instance_id: string;
    started_at: string;
    finished_at: string;
    duration_seconds: number;
    status: Crawl.status;
    number_of_peers?: number;
    total_users?: number;
    active_users_half_year?: number;
    active_users_month?: number;
    local_posts?: number;
    local_comments?: number;
    raw_nodeinfo?: Record<string, any>;
};

export namespace Crawl {

    export enum status {
        COMPLETED = 'completed',
        FAILED = 'failed',
        BLOCKED = 'blocked',
        TIMEOUT = 'timeout',
        INTERNAL_ERROR = 'internal_error',
    }


}

