import LZString from 'lz-string';

/**
 * Encodes profile object to a compressed base64 string
 * @param {Object} data 
 * @returns {string}
 */
export const encodeProfile = (data) => {
    try {
        const jsonStr = JSON.stringify(data);
        return LZString.compressToEncodedURIComponent(jsonStr);
    } catch (error) {
        console.error('Encoding error:', error);
        return '';
    }
};

/**
 * Decodes compressed base64 string to profile object
 * @param {string} encodedStr 
 * @returns {Object|null}
 */
export const decodeProfile = (encodedStr) => {
    try {
        const jsonStr = LZString.decompressFromEncodedURIComponent(encodedStr);
        if (!jsonStr) return null;
        return JSON.parse(jsonStr);
    } catch (error) {
        console.error('Decoding error:', error);
        return null;
    }
};
