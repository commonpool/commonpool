import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {ActivatedRoute, Router} from '@angular/router';
import {distinctUntilChanged, filter, pluck, switchMap, tap} from 'rxjs/operators';
import {AuthService} from '../../auth.service';
import {Observable} from 'rxjs';
import {ExtendedResource} from '../../api/models';

@Component({
  selector: 'app-resource-details',
  templateUrl: './resource-details.component.html',
  styleUrls: ['./resource-details.component.css']
})
export class ResourceDetailsComponent implements OnInit {

  resourceId$ = this.route.params.pipe(
    pluck('id'),
    filter(r => !!r),
    distinctUntilChanged());
  resource$: Observable<ExtendedResource> = this.resourceId$.pipe(
    tap(console.log),
    switchMap(id => this.backend.getResource(id)),
    pluck('resource')
  );

  constructor(
    private backend: BackendService,
    private router: Router,
    private route: ActivatedRoute,
    public auth: AuthService
  ) {}

  async editResource(id: string) {
    await this.router.navigateByUrl('/resources/' + id + '/edit');
  }

  ngOnInit(): void {
  }

}
