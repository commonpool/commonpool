import {Injectable} from '@angular/core';
import {Subject} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ResourceEvaluationService {

  private resourceEvaluationSubject = new Subject<string>();
  resourceEvaluation$ = this.resourceEvaluationSubject.asObservable();

  evaluateResource(resourceId: string) {
    this.resourceEvaluationSubject.next(resourceId);
  }

}
