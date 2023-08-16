/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

export type Server = {
    id: string;
    domain: string;
    status: Server.status;
    description?: string;
    software?: string;
    version?: string;
    number_of_peers?: number;
    open_registrations?: boolean;
    total_users?: number;
    active_users_half_year?: number;
    active_users_month?: number;
    local_posts?: number;
    local_comments?: number;
};

export namespace Server {

    export enum status {
        ACTIVE = 'active',
        INACTIVE = 'inactive',
        UNKNOWN = 'unknown',
    }


}

