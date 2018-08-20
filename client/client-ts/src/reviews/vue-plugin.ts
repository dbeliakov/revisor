import _Vue from 'vue';
import ReviewsService from '@/reviews/service';
import {AxiosStatic} from 'axios';

export function ReviewsPlugin(Vue: typeof _Vue, options: {axios: AxiosStatic}): void {
    const service = new ReviewsService(options.axios);
    Vue.prototype.$reviews = service;
}
