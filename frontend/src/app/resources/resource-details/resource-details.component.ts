import {Component, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {ActivatedRoute, Router} from '@angular/router';
import {distinctUntilChanged, pluck, switchMap} from 'rxjs/operators';
import {AuthService} from '../../auth.service';

@Component({
  selector: 'app-resource-details',
  templateUrl: './resource-details.component.html',
  styleUrls: ['./resource-details.component.css']
})
export class ResourceDetailsComponent implements OnInit {

  resourceId$ = this.route.params.pipe(pluck('id'), distinctUntilChanged());
  resource$ = this.resourceId$.pipe(
    switchMap(id => this.backend.getResource(id)),
    pluck('resource')
  );

  constructor(
    private backend: BackendService,
    private router: Router,
    private route: ActivatedRoute,
    public auth: AuthService
  ) {

  }

  async editResource(id: string) {
    await this.router.navigateByUrl('/resources/' + id + '/edit');
  }

  ngOnInit(): void {
  }

}
