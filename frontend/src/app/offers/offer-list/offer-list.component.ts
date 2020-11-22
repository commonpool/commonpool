import {Component, Input, OnInit} from '@angular/core';
import {BackendService} from '../../api/backend.service';
import {AcceptOfferRequest, DeclineOfferRequest, GetOffersRequest, Offer} from '../../api/models';
import {ActivatedRoute} from '@angular/router';
import {AuthService} from '../../auth.service';
import {distinctUntilChanged, filter, pluck, shareReplay, switchMap, tap} from 'rxjs/operators';
import {ReplaySubject, Subject} from 'rxjs';

@Component({
  selector: 'app-offer-list',
  templateUrl: './offer-list.component.html',
  styleUrls: ['./offer-list.component.css']
})
export class OfferListComponent implements OnInit {

  constructor(private backend: BackendService, private route: ActivatedRoute, private auth: AuthService) {

  }

  userId$ = this.auth.authUserId$.pipe(
    filter(uid => !!uid),
    distinctUntilChanged(),
    tap((uid) => this.refresh()),
    shareReplay()
  );

  refreshSubject = new ReplaySubject();
  offers$ = this.refreshSubject.asObservable().pipe(
    switchMap(() => this.backend.getOffers(new GetOffersRequest())),
    pluck('offers'),
    shareReplay()
  );

  ngOnInit(): void {
    this.refresh();
  }

  accept(id: string) {
    this.backend.acceptOffer(new AcceptOfferRequest(id)).subscribe(res => {
      this.refreshSubject.next();
    });
  }

  decline(id: string) {
    this.backend.declineOffer(new DeclineOfferRequest(id)).subscribe(res => {
      this.refreshSubject.next();
    });
  }

  refresh() {
    this.refreshSubject.next();
  }

}
