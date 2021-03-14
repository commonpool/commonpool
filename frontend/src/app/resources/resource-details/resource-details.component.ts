import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {ActivatedRoute, Router} from '@angular/router';
import {delay, distinctUntilChanged, filter, pluck, shareReplay, startWith, switchMap, tap} from 'rxjs/operators';
import {AuthService} from '../../auth.service';
import {combineLatest, Observable, of, Subject} from 'rxjs';
import {ExtendedResource} from '../../api/models';
import {ResourceEvaluationService} from '../resource-evaluation/resource-evaluation.service';

@Component({
  selector: 'app-resource-details',
  templateUrl: './resource-details.component.html',
  styleUrls: ['./resource-details.component.css']
})
export class ResourceDetailsComponent implements OnInit {

  refreshSubject = new Subject();
  refresh$ = this.refreshSubject.pipe(delay(200), startWith(true));

  resourceId$ = this.route.params.pipe(
    pluck('id'),
    filter(r => !!r),
    shareReplay());

  resource$: Observable<ExtendedResource> = combineLatest([this.resourceId$, this.refresh$]).pipe(
    switchMap(([id]) => this.backend.getResource(id)),
    pluck('resource')
  );

  constructor(
    private backend: BackendService,
    private router: Router,
    private route: ActivatedRoute,
    public auth: AuthService,
    public resourceEvaluationService: ResourceEvaluationService
  ) {
  }

  async editResource(id: string) {
    await this.router.navigateByUrl('/resources/' + id + '/edit');
  }

  public evaluateResource(resourceId: string) {
    this.resourceEvaluationService.evaluateResource(of(resourceId)).subscribe((ok) => {
      if (ok) {
        this.refreshSubject.next(true);
      }
    });
  }

  ngOnInit(): void {
  }

}
