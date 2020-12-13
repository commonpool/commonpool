import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {
  ConfirmBorrowedResourceReturned,
  ConfirmResourceBorrowed,
  ConfirmResourceTransferred,
  ConfirmServiceProvidedRequest,
  Offer
} from '../../api/models';
import {combineLatest, of, ReplaySubject, Subject} from 'rxjs';
import {map} from 'rxjs/operators';
import {AuthService} from '../../auth.service';
import {BackendService} from '../../api/backend.service';

@Component({
  selector: 'app-offer-list-item',
  templateUrl: './offer-list-item.component.html',
  styleUrls: ['./offer-list-item.component.css']
})
export class OfferListItemComponent implements OnInit {

  constructor(public auth: AuthService, private backend: BackendService) {
  }

  private offerSubject = new ReplaySubject<Offer>(1);
  public offerSubject$ = this.offerSubject.asObservable();

  public offerItems$ = this.offerSubject$.pipe(
    map(o => o.items)
  );

  public offerItemsICanApprove$ = combineLatest([this.offerItems$, this.auth.authUserId$]).pipe(
    map(([offerItems, userId]) => {
      return offerItems.filter(offerItem => {
        return (offerItem.givingApprovers.includes(userId) || offerItem.receivingApprovers.includes(userId));
      });
    })
  );

  @Input()
  set offer(value: Offer) {
    this.offerSubject.next(value);
  }

  @Output()
  approve = new EventEmitter();

  @Output()
  decline = new EventEmitter();

  @Output()
  refresh = new EventEmitter();

  confirmServiceProvided(id: string) {
    const sub = this.backend.confirmServiceProvided(new ConfirmServiceProvidedRequest(id)).subscribe(res => {
      sub.unsubscribe();
      this.refresh.next();
    });
  }

  confirmResourceTransferred(id: string) {
    const sub = this.backend.confirmResourceTransfer(new ConfirmResourceTransferred(id)).subscribe(res => {
      sub.unsubscribe();
      this.refresh.next();
    });
  }

  confirmResourceBorrowed(id: string) {
    const sub = this.backend.confirmResourceBorrowed(new ConfirmResourceBorrowed(id)).subscribe(res => {
      sub.unsubscribe();
      this.refresh.next();
    });
  }

  confirmResourceBorrowedReturned(id: string) {
    const sub = this.backend.confirmBorrowedResourceReturned(new ConfirmBorrowedResourceReturned(id)).subscribe(res => {
      sub.unsubscribe();
      this.refresh.next();
    });
  }

  ngOnInit(): void {
  }

}
