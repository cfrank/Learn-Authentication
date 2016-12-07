/* @flow */
export type AuthenticationData = {
    authString: string,
    date: number,
    nonce: string
};

export const NONCELEN: number = 12;

export async function authCall(data: AuthenticationData, url: string): Promise<string>{
    try{
        let response: Response = await fetch(url, {
            // Make the request to the api
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        let responseText: string = await response.json();
        return responseText;
    }catch(e){
        throw new Error("Error making authentication request");
    }
}

export function generateNonce(length: number): string{
    let keySpace: string = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_';
    let nonce: string = '';
    for(let i: number = 0; i <= length; ++i){
        nonce += keySpace.charAt(Math.floor(Math.random() * keySpace.length));
    }
    return nonce;
}

function _checkResponse(){
    console.log('yo');
}