import {Injectable} from '@angular/core';
import {ErrorResponse, SearchResourceRequest, SearchResourcesResponse} from './models';
import {Observable, throwError} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {catchError, map, tap} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class BackendService {

  constructor(private http: HttpClient) {

  }

  searchResources(request: SearchResourceRequest): Observable<SearchResourcesResponse> {
    return this.http.get<SearchResourcesResponse>(`http://127.0.0.1:8585/api/v1/resources`, {
      params: {
        take: request.take.toString(),
        skip: request.skip.toString(),
        query: request.query,
        type: request.type.toString(),
      },
      observe: 'response'
    }).pipe(
      tap(console.log),
      map((res) => {
        if (res.status > 399 || res.body === null) {
          console.log('Error');
          throwError(ErrorResponse.fromHttpResponse(res));
        }
        const body = res.body as SearchResourcesResponse;
        console.log(body)
        return new SearchResourcesResponse(body.resources, body.totalCount, body.take, body.skip);
      })
    );
  }

}
