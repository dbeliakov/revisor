import { AxiosStatic } from '../../node_modules/axios';
import Review from '@/reviews/review';
import { responseToError } from '@/utils/utils';
import {UserInfo} from '@/auth/user-info';
import { Diff } from '@/reviews/diff';
import Comment from '@/reviews/comment';

export class DiffReply {
    public info: Review;
    public diff: Diff;
    public comments: Comment[];

    public constructor(json: any) {
        this.info = new Review(json.info);
        this.diff = new Diff(json.diff);
        this.comments = [];
        for (const comment of json.comments) {
            this.comments.push(new Comment(comment));
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

    public async updateReview(
            reviewId: number,
            name: string,
            reviewers: string,
            fileName: string | null,
            fileContent: string | null): Promise<Error | undefined> {
        try {
            const params: {[id: string]: string; } = {name, reviewers};
            if (fileName !== null && fileContent !== null) {
                params.new_revision = fileContent;
            }
            await this.axios.post('/reviews/' + reviewId + '/update', params);
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
            reviewId: number,
            startRevision: number | null,
            endRevision: number | null): Promise<DiffReply | Error> {
        let path = '/reviews/' + reviewId;
        if (startRevision && endRevision) {
            path += '?start_rev=' + (startRevision - 1) + '&end_rev=' + (endRevision! - 1);
        }
        try {
            const response = await this.axios.get(path);
            return new DiffReply(response.data.data);
        } catch (error) {
            return responseToError(error);
        }
    }

    public async addComment(lineId: string, reviewId: number, text: string, parentId: number = 0)
            : Promise<Error | undefined> {
        try {
            const data = {
                review_id: Number(reviewId),
                line_id: lineId,
                text,
            };
            if (parentId > 0) {
                (data as any).parent = parentId;
            }
            await this.axios.post('/comments/add', data);
        } catch (error) {
            return responseToError(error);
        }
    }

    public async acceptReview(reviewId: number): Promise<Error | undefined> {
        try {
            await this.axios.post('/reviews/' + reviewId + '/accept');
        } catch (error) {
            return responseToError(error);
        }
    }

    public async declineReview(reviewId: number): Promise<Error | undefined> {
        try {
            await this.axios.post('/reviews/' + reviewId + '/decline');
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
