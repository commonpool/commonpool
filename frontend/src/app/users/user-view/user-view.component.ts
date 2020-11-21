import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {pluck} from 'rxjs/operators';
import {UserInfoService} from '../user-info.service';
import {Observable, Subscription} from 'rxjs';
import {UserInfoResponse} from '../../api/models';

@Component({
  selector: 'app-user-view',
  templateUrl: './user-view.component.html',
  styleUrls: ['./user-view.component.css']
})
export class UserViewComponent implements OnInit, OnDestroy {

  userSubscription: Subscription;
  user$: Observable<UserInfoResponse>;

  constructor(private route: ActivatedRoute, public userService: UserInfoService) {
    this.user$ = this.userService.getUserInfo();
  }

  ngOnInit(): void {
    this.userSubscription = this.route.params.pipe(
      pluck('id')
    ).subscribe(id => {
      this.userService.setUserId(id);
    });
  }

  ngOnDestroy() {
    if (this.userSubscription) {
      this.userSubscription.unsubscribe();
    }
  }

}
