import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {BackendService} from '../../api/backend.service';
import {pluck, switchMap, tap} from 'rxjs/operators';

@Component({
  selector: 'app-user-view',
  templateUrl: './user-view.component.html',
  styleUrls: ['./user-view.component.css']
})
export class UserViewComponent implements OnInit {

  constructor(private route: ActivatedRoute, private backend: BackendService) {

  }

  userId$ = this.route.params.pipe(pluck('id'));
  user$ = this.userId$.pipe(
    switchMap(id => this.backend.getUserInfo(id))
  );

  ngOnInit(): void {
  }

}
