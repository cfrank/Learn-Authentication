/* @flow */
import {authCall, generateNonce, NONCELEN} from './authentication';
import type {AuthenticationData} from './authentication';

type AuthenticationInput = {
    container: HTMLElement,
    input: HTMLInputElement
};

export class AuthenticationForm{
    authForm: HTMLFormElement; // The form element
    authInputs: Array<AuthenticationInput>; // Object with fields and inputs

    constructor(authForm: HTMLFormElement){
        // Initialize variables
        this.authForm = authForm;
        this.authInputs = [];
        // Call methods
        this._getElements();
    }

    /*
     * This is a generic validation for a authentication form
     * any more thorough validation must be done by the child class
     *
     * @returns Array containing errors
     */ 
    validate(): boolean{
        let errorList: Array<number> = [];
        for(let i: number = 0; i < this.authInputs.length; ++i){
            let value: string = this.authInputs[i].input.value,
                inputType: string = this.authInputs[i].input.getAttribute('type');
            switch(inputType){
                case 'email':
                    let emailRegex: RegExp = /\S+@\S+\.\S+/;
                    if(emailRegex.test(value) !== true){
                        errorList.push(i);
                    }
                    break;
                default:
                    if(!value || value.length <= 0){
                        console.log('no val' + i);
                        errorList.push(i);
                    }
                    break;
            }
        }
        this._showErrors(errorList);

        return (errorList.length > 0) ? false : true;
    }

    /*
     * Base submit method, anything more thorough must be overridden by
     * child class
     */
    submit(): void{
        if(this.validate()){
            // Disable submit
            this.toggleSubmit(false);
        }
        else{
            // Enable submit
            this.toggleSubmit(true);
        }
    }

    /*
     * Toggle the submit button so that a double click on submit
     * doesn't cause two submissions
     *
     * @param Boolean - Submit enabled?
     */
    toggleSubmit(enabled: boolean): void{
        if(enabled){
            this.authForm.submit.disabled = false;
        }
        else if(!this.authForm.submit.disabled && !enabled){
            this.authForm.submit.disabled = true;
        }
    }

    /*
     * Show errors for the form
     *
     * @param Array<number> - Indexes of inputs with errors
     */
    _showErrors(errorList: Array<number>): void{
        for(let i: number = 0; i < this.authInputs.length; ++i){
            // Check if the input has an error
            if(errorList.indexOf(i) >= 0){
                this._toggleError(i, true);
            }
            else{
                this._toggleError(i, false);
            }
        }
    }

    /*
     * Populate values with their fields
     */
    _getElements(): void{
        let fields: NodeList<HTMLElement> = this.authForm.querySelectorAll('div.form-field');
        for(let i: number = 0; i < fields.length; ++i){
            let input: AuthenticationInput = {
                container: fields[i],
                input: fields[i].getElementsByTagName('input')[0]
            }
            this.authInputs.push(input);
        }
    }

    /*
     * Check if an element contains the error class and remove it
     *
     * @param Number - The index to the element
     * @param Boolean - Showing the error? Or removing it.
     */
    _toggleError(index: number, show: boolean): void{
        let element: HTMLElement = this.authInputs[index].container,
            hasClass: boolean = element.classList.contains('error');
        if(show && !hasClass){
            element.classList.add('error');
        }
        else if(!show && hasClass){
            element.classList.remove('error');
        }
    }
}

export class SignIn extends AuthenticationForm{
    constructor(authForm: HTMLFormElement){
        super(authForm);
    }

    submit(): void{
        if(!super.validate()){
            super.toggleSubmit(true);
            console.log('ERROR');
        }
        else{
            super.toggleSubmit(false);
            let data: AuthenticationData = {
                authString: this._getAuthString(),
                date: ~~(Date.now() / 1000),
                nonce: generateNonce(NONCELEN)
            };

            console.log(data);
        }
    }

    _getAuthString(): string{
        let username: string = this.authInputs[0].input.value;
        let password: string = this.authInputs[1].input.value;
        let authString: string = `${username}&${password}`;
        let encodedString: string = window.btoa(authString);
        return encodedString;
    }
}

export class ForgotPassword extends AuthenticationForm{
    constructor(authForm: HTMLFormElement){
        super(authForm);
    }
}

export class SignUp extends AuthenticationForm{
    constructor(authForm: HTMLFormElement){
        super(authForm);
    }

    submit(): void{
        if(!this._checkPasswords() || !super.validate()){
            super.toggleSubmit(true);
            console.log('ERROR');
        }
        else{
            super.toggleSubmit(false);
            let data: AuthenticationData = {
                authString: this._getAuthString(),
                date: ~~(Date.now() / 1000),
                nonce: generateNonce(NONCELEN)
            };
            
            let hello = authCall(data, '/auth/signup');
            console.log(hello);
        }
    }

    _checkPasswords(): boolean{
        let password: string = this.authInputs[1].input.value;
        let confirm: string = this.authInputs[2].input.value;

        if(password !== confirm){
            this.authInputs[2].container.classList.add('error');
            return false;
        }
        return true;
    }

    _getAuthString(): string{
        let email: string = this.authInputs[0].input.value;
        let password: string = this.authInputs[1].input.value;
        let authString: string = `${email}&${password}`
        let encodedString: string = window.btoa(authString);
        return encodedString;
    }
}