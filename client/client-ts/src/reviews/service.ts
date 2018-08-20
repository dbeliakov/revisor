import { AxiosStatic } from '../../node_modules/axios';
import Review from '@/reviews/review';
import { responseToError } from '@/utils/utils';
import {UserInfo} from '@/auth/user-info';
import { Diff } from '@/reviews/diff';
import Comment from '@/reviews/comment';

class DiffReply {
    public info: Review;
    public diff: Diff;
    public comments: Comment[];

    public constructor(json: any) {
        this.info = new Review(json.info);
        this.diff = new Diff(json.diff);
        this.comments = [];
        for (const comment of json.comments) {
            this.comments.push(new Comment(json));
        }
    }
}

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

    public async searchReviewers(query: string): Promise<UserInfo[] | Error> {
        try {
            const response = await this.axios.get('/users/search?query=' + query);
            const result: UserInfo[] = [];
            for (const user of response.data.data) {
                result.push(new UserInfo(user));
            }
            return result;
        } catch (error) {
            return responseToError(error);
        }
    }

    public async loadDiff(
            reviewId: string,
            startRevision: number | undefined,
            endRevision: number | undefined): Promise<DiffReply | Error> {
        let path = '/reviews/' + reviewId;
        if (startRevision && endRevision) {
            path += '?start_rev=' + startRevision + '&end_rev=' + endRevision;
        }
        try {
            const response = await this.axios.get(path);
            return new DiffReply(response);
        } catch (error) {
            return responseToError(error);
        }
    }

    public async addComment(lineId: string, reviewId: string, text: string, parentId: string = '')
            : Promise<Error | undefined> {
        try {
            await this.axios.post('/comments/add', {
                review_id: reviewId,
                line_id: lineId,
                text,
                parent: parentId,
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
