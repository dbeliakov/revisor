import _Vue from 'vue';
import Auth from '@/auth/auth';
import {AxiosStatic} from 'axios';
import Router from 'vue-router';

export function AuthPlugin(Vue: typeof _Vue, options: {axios: AxiosStatic, router: Router}): void {
    const auth = new Auth(options.axios);
    Vue.prototype.$auth = auth;
    options.router.beforeEach((to, from, next) => {
        auth.whenBecomeReady(() => {
            if (to.matched.some((record) => record.meta.requiresAuth)) {
                if (auth.authenticated()) {
                    next();
                } else {
                    next({
                        name: 'login',
                    });
                }
            } else if (to.matched.some((record) => record.meta.requiresNoAuth)) {
                if (!auth.authenticated()) {
                    next();
                } else {
                    next({
                        name: 'home',
                    });
                }
            } else {
                next();
            }
        });
    });
}
