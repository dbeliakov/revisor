import { AxiosStatic } from '../../node_modules/axios';
import Review from '@/reviews/review';
import { responseToError } from '@/utils/utils';

export default class ReviewsService {
    private axios: AxiosStatic;

    public constructor(axios: AxiosStatic) {
        this.axios = axios;
    }

    public async loadIncomingReviews(): Promise<Review[] | Error> {
        return this.loadReviews('/reviews/incoming');
    }

    public async loadOutgoingReviews(): Promise<Review[] | Error> {
        return this.loadReviews('/reviews/outgoing');
    }

    public async createReview(
            name: string,
            reviewers: string,
            fileName: string,
            fileContent: string): Promise<Error | undefined> {
        try {
            await this.axios.post('/reviews/new', {
                name,
                reviewers,
                file_content: fileContent,
                file_name: fileName,
            });
        } catch (error) {
            return responseToError(error);
        }
    }

    private async loadReviews(path: string): Promise<Review[] | Error> {
        try {
            const response = await this.axios.get(path);
            const result: Review[] = [];
            for (const review of response.data.data) {
                result.push(new Review(review));
            }
            return result;
        } catch (error) {
            return responseToError(error);
        }
    }
}
