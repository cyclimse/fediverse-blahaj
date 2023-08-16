/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Crawl } from '../models/Crawl';
import type { Error } from '../models/Error';
import type { Server } from '../models/Server';

import type { CancelablePromise } from '../core/CancelablePromise';
import type { BaseHttpRequest } from '../core/BaseHttpRequest';

export class DefaultService {

    constructor(public readonly httpRequest: BaseHttpRequest) {}

    /**
     * List all servers
     * @param software filter by software name.
     * @param page page number of results to return
     * @param perPage number of results to return per page
     * @returns any paginated array of servers
     * @returns Error unexpected error
     * @throws ApiError
     */
    public listServers(
        software?: string,
        page: number = 1,
        perPage: number = 30,
    ): CancelablePromise<{
        results: Array<Server>;
        total: number;
        page: number;
        per_page: number;
    } | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/servers',
            query: {
                'software': software,
                'page': page,
                'per_page': perPage,
            },
        });
    }

    /**
     * Info for a specific server
     * @param id ID of the server to fetch
     * @returns Server server response
     * @returns Error unexpected error
     * @throws ApiError
     */
    public getServerById(
        id: string,
    ): CancelablePromise<Server | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/servers/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * List all crawls for a server
     * @param id ID of the server to fetch
     * @param page page number of results to return
     * @param perPage number of results to return per page
     * @returns any paginated array of crawls
     * @returns Error unexpected error
     * @throws ApiError
     */
    public listCrawlsForServer(
        id: string,
        page: number = 1,
        perPage: number = 30,
    ): CancelablePromise<{
        results: Array<Crawl>;
        total: number;
        page: number;
        per_page: number;
    } | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/servers/{id}/crawls',
            path: {
                'id': id,
            },
            query: {
                'page': page,
                'per_page': perPage,
            },
        });
    }

}
