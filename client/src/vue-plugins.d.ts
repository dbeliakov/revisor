import Auth from '@/auth/auth';
import ReviewsService from '@/reviews/service';

declare module 'vue/types/vue' {
  interface Vue {
    $auth: Auth;
    $reviews: ReviewsService;
  }
}
