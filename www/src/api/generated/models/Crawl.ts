/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Crawl = {
    id: string;
    server_id: string;
    started_at: string;
    finished_at: string;
    status: Crawl.status;
    error?: string;
    number_of_peers?: number;
    total_users?: number;
    active_users_half_year?: number;
    active_users_month?: number;
    local_posts?: number;
    local_comments?: number;
};

export namespace Crawl {

    export enum status {
        PENDING = 'pending',
        RUNNING = 'running',
        FINISHED = 'finished',
        FAILED = 'failed',
    }


}

