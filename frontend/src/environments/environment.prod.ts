/* tslint:disable:no-string-literal */
import {IEnvironment} from './IEnvironment';

export const environment: IEnvironment = {
  production: true,
  apiUrl: window['env']['apiUrl'] || 'default',
  debug: window['env']['debug'] || false
};
