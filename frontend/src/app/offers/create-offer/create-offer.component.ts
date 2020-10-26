import {Component, OnInit} from '@angular/core';
import {OfferItemType, SendOfferRequest, SendOfferRequestItem, SendOfferRequestPayload} from '../../api/models';
import {BackendService} from '../../api/backend.service';

@Component({
  selector: 'app-create-offer',
  templateUrl: './create-offer.component.html',
  styleUrls: ['./create-offer.component.css']
})
export class CreateOfferComponent implements OnInit {

  constructor(private backend: BackendService) {
  }

  pending = false;
  error = undefined;
  items: SendOfferRequestItem[] = [];
  newFrom: string | null = null;
  newTo: string | null = null;
  resource: string | null = null;
  offerItemType: OfferItemType = OfferItemType.ResourceItem;
  time = 1;

  toPredicate = (val: string) => true;

  setNewFrom(value: string) {
    this.newFrom = value;
    this.toPredicate = (val: string) => val !== value;
  }

  ngOnInit(): void {
    this.items = [];
  }

  add() {
    this.items.push(new SendOfferRequestItem(this.newFrom, this.newTo, this.offerItemType, this.resource, this.time * 60 * 60));
    this.newFrom = null;
    this.newTo = null;
    this.resource = null;
    this.time = 1;
  }


  remove(i: number) {
    this.items.splice(i, 1);
    this.items = [...this.items];
  }

  submit() {
    this.pending = true;
    this.error = undefined;
    this.backend.sendOffer(new SendOfferRequest(new SendOfferRequestPayload(this.items))).subscribe(res => {
      this.pending = false;
      console.log(res);
    }, err => {
      this.pending = false;
      this.error = err;
    });
  }

}
