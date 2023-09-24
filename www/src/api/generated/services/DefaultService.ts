/* generated using openapi-typescript-codegen -- do no edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Crawl } from '../models/Crawl';
import type { Error } from '../models/Error';
import type { Instance } from '../models/Instance';

import type { CancelablePromise } from '../core/CancelablePromise';
import type { BaseHttpRequest } from '../core/BaseHttpRequest';

export class DefaultService {

    constructor(public readonly httpRequest: BaseHttpRequest) {}

    /**
     * List all instances
     * @param software filter by software name.
     * @param page page number of results to return
     * @param perPage number of results to return per page
     * @returns any paginated array of instances
     * @returns Error unexpected error
     * @throws ApiError
     */
    public listInstances(
        software?: string,
        page: number = 1,
        perPage: number = 30,
    ): CancelablePromise<{
        results: Array<Instance>;
        total: number;
        page: number;
        per_page: number;
    } | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/instances',
            query: {
                'software': software,
                'page': page,
                'per_page': perPage,
            },
        });
    }

    /**
     * Info for a specific instance
     * @param id ID of the instance to fetch
     * @returns Instance instance response
     * @returns Error unexpected error
     * @throws ApiError
     */
    public getInstanceById(
        id: string,
    ): CancelablePromise<Instance | Error> {
        return this.httpRequest.request({
            method: 'GET',
            url: '/instances/{id}',
            path: {
                'id': id,
            },
        });
    }

    /**
     * List all crawls for a instance
     * @param id ID of the instance to fetch
     * @param page page number of results to return
     * @param perPage number of results to return per page
     * @returns any paginated array of crawls
     * @returns Error unexpected error
     * @throws ApiError
     */
    public listCrawlsForInstance(
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
            url: '/instances/{id}/crawls',
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
